package chatgpt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"notionboy/internal/pkg/logger"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"

	"github.com/donovanhide/eventsource"
)

type reverseClient struct {
	mu           sync.Mutex
	SessionToken string `json:"session_token"`
	AuthToken    string `json:"auth_token"`
	isRateLimit  atomic.Bool
}

func newReverseClient(sessionToken string) Chatter {
	once := sync.Once{}
	var client *reverseClient
	once.Do(func() {
		client = &reverseClient{
			SessionToken: sessionToken,
		}
		client.setIsRateLimit(false)
		client.refreshSession()
		go func() {
			for range time.Tick(5 * time.Minute) {
				client.refreshSession()
			}
		}()
	})
	return client
}

func (cli *reverseClient) Chat(ctx context.Context, parentMessageId, prompt string) (string, string, error) {
	if cli.GetIsRateLimit() {
		return "", "", errors.New("hit rate limit, please retry after 1 hour")
	}

	cli.mu.Lock()
	defer cli.mu.Unlock()
	bodyChan := make(chan *http.Response)
	defer close(bodyChan)
	errorChan := make(chan error)
	defer close(errorChan)
	stopRetryChan := make(chan struct{})
	payload := buildPayload(parentMessageId, []string{prompt})
	logger.SugaredLogger.Debugw("Request Payload", "payload", payload)

	post := func() {
		resp, err := client.R().
			SetContext(ctx).
			SetBody(payload).
			SetHeader("Accept", "text/event-stream").
			SetDoNotParseResponse(true).
			Post(chatURL)
		if err != nil {
			errorChan <- err
			return
		}

		if resp.StatusCode() == http.StatusBadRequest {
			cli.setIsRateLimit(true)
			logger.SugaredLogger.Errorw("Reach ratelimit, please try after 1 hour", "status", resp.Status())
			stopRetryChan <- struct{}{}
			return
		}
		logger.SugaredLogger.Debugw("status", "status_code", resp.Status())
		if resp.StatusCode() != http.StatusOK {
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
			res, err = cli.decodeResponse(body)
			if err != nil {
				err = fmt.Errorf("Decode response error, %s", err.Error())
				continue
			}
			return res.Message.ID, res.Message.Content.Parts[0], err
		case err = <-errorChan:
			logger.SugaredLogger.Warnw("Get response from chatGPT error", "retry_times", i+1, "err", err)
			continue
		case <-stopRetryChan:
			break
		}
	}
	return "", "", err
}

func (cli *reverseClient) decodeResponse(resp *http.Response) (*ResponseBody, error) {
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

func (cli *reverseClient) GetIsRateLimit() bool {
	return cli.isRateLimit.Load()
}

func (cli *reverseClient) setIsRateLimit(flag bool) {
	cli.isRateLimit.Store(flag)
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
