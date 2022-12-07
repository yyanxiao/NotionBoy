package chatgpt

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"notionboy/internal/pkg/logger"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/donovanhide/eventsource"
)

type chatter struct {
	mu           sync.Mutex
	SessionToken string `json:"session_token"`
	AuthToken    string `json:"auth_token"`
}

type Chatter interface {
	Chat(parentMessageId, prompt string) (string, string, error)
}

var bot *chatter

// New create a new chatter
func New(sessionToken string) Chatter {
	refreshSession()
	once := sync.Once{}
	once.Do(func() {
		go func() {
			for range time.Tick(5 * time.Minute) {
				refreshSession()
			}
		}()
		bot = &chatter{
			SessionToken: sessionToken,
		}
	})
	return bot
}

func (c *chatter) Chat(parentMessageId, prompt string) (string, string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	bodyChan := make(chan *http.Response)
	defer close(bodyChan)
	errorChan := make(chan error)
	defer close(errorChan)
	payload := buildPayload(parentMessageId, []string{prompt})
	logger.SugaredLogger.Debugw("Request Payload", "payload", payload)

	post := func() {
		resp, err := client.R().
			SetBody(payload).
			SetHeader("Accept", "text/event-stream").
			SetDoNotParseResponse(true).
			Post(chatURL)
		if err != nil {
			errorChan <- err
			return
		}

		logger.SugaredLogger.Debugw("status", "status_code", resp.Status())
		if resp.StatusCode() != http.StatusOK {
			refreshSession()
			errorChan <- errors.New("Status: " + resp.Status())
			return
		}
		bodyChan <- resp.RawResponse
	}

	var err error
	// if chatGPT response with error, retry 3 times
	for i := 0; i < 3; i++ {
		go post()
		select {
		case body := <-bodyChan:
			var res *ResponseBody
			res, err = c.decodeResponse(body)
			if err != nil {
				err = fmt.Errorf("Decode response error, %s", err.Error())
				continue
			}
			return res.Message.ID, res.Message.Content.Parts[0], err
		case err = <-errorChan:
			logger.SugaredLogger.Warnw("Get response from chatGPT error", "retry_times", i+1, "err", err)
			continue
		}
	}
	return "", "", err
}

func (c *chatter) decodeResponse(resp *http.Response) (*ResponseBody, error) {
	defer resp.Body.Close()
	eventChan := make(chan string)

	decoder := eventsource.NewDecoder(resp.Body)
	go func() {
		defer close(eventChan)
		for {
			event, err := decoder.Decode()
			if err != nil {
				logger.SugaredLogger.Errorw("Failed to decode event", "err", err)
				break
			}
			if event.Data() == "[DONE]" || event.Data() == "" {
				break
			}
			// logger.SugaredLogger.Debugw("send chunk data", "data", string(event.Data()))
			eventChan <- event.Data()
		}
	}()

	var res ResponseBody
	for chunk := range eventChan {
		if err := json.Unmarshal([]byte(chunk), &res); err != nil {
			continue
		}
		//if len(res.Message.Content.Parts) > 0 {
		//	logger.SugaredLogger.Debug(res.Message.Content.Parts[0])
		//}
	}

	if len(res.Message.Content.Parts) == 0 {
		return nil, errors.New("ChatGPT do not response, please retry later")
	}
	logger.SugaredLogger.Debugw("Response", "conversation_id", res.ConversationId, "error", res.Error, "message", res.Message)
	return &res, nil
}

func buildPayload(parentMessageId string, prompt []string) *Payload {
	if parentMessageId == "" {
		parentMessageId = generateUUID()
	}
	return &Payload{
		Action: "next",
		Messages: []Message{
			{
				ID:   generateUUID(),
				Role: "user",
				Content: MessageContent{
					ContentType: "text",
					Parts:       prompt,
				},
			},
		},
		// ConversationId:  conversationId,
		ParentMessageId: parentMessageId,
		Model:           "text-davinci-002-render",
	}
}

func generateUUID() string {
	uid := uuid.New()
	return uid.String()
}
