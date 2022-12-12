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

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"

	"github.com/donovanhide/eventsource"
)

const (
	ERROR_MSG_INTERNAL_ERROR = iota + 1
	ERROR_MSG_RATE_LIMITED_ERROR
	ERROR_MSG_AUTH_ERROR
	ERROR_MSG_UNKNOWN_ERROR
)

type errorResp struct {
	Type   int
	Status string
}

type reverseClient struct {
	mu           sync.Mutex
	client       *resty.Client
	SessionToken string `json:"session_token"`
	authToken    string
	Email        string `json:"email"`
	Password     string `json:"password"`
	isRateLimit  atomic.Bool
}

func newReverseClient(sessionToken, email, password string) Chatter {
	client := &reverseClient{
		SessionToken: sessionToken,
		Email:        email,
		Password:     password,
		client:       resty.New(),
	}
	client.setIsRateLimit(false)
	client.refreshSession()

	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			client.refreshSession()
		}
	}()
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
	errorChan := make(chan *errorResp)
	defer close(errorChan)

	payload := buildPayload(parentMessageId, []string{prompt})
	logger.SugaredLogger.Debugw("Request Payload", "payload", payload)

	var err error
	// if chatGPT response with error, retry 3 times
L:
	for i := 0; i < 3; i++ {
		go cli.sendRequest(ctx, payload, bodyChan, errorChan)
		select {
		case body := <-bodyChan:
			var res *ResponseBody
			res, err = cli.decodeResponse(body)
			if err != nil {
				err = fmt.Errorf("Decode response error, %s", err.Error())
				continue L
			}
			return res.Message.ID, res.Message.Content.Parts[0], err
		case errResp := <-errorChan:
			switch errResp.Type {
			case ERROR_MSG_UNKNOWN_ERROR, ERROR_MSG_INTERNAL_ERROR:
				logger.SugaredLogger.Warnw("Get response from chatGPT error", "retry_times", i+1, "status", errResp.Status)
			case ERROR_MSG_AUTH_ERROR, ERROR_MSG_RATE_LIMITED_ERROR:
				logger.SugaredLogger.Warnw("ChatGPT auth error, exit", "retry_times", i+1, "status", errResp.Status)
				break L
			}
		}
	}
	return "", "", err
}

func (cli *reverseClient) sendRequest(ctx context.Context, payload *Payload, bodyChan chan *http.Response, errorChan chan *errorResp) {
	resp, err := cli.client.R().
		SetContext(ctx).
		SetBody(payload).
		SetHeader("Accept", "text/event-stream").
		SetHeader("Authorization", "Bearer "+cli.authToken).
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", userAgent).
		SetDoNotParseResponse(true).
		Post(chatURL)
	if err != nil {
		errorChan <- &errorResp{
			Type: ERROR_MSG_UNKNOWN_ERROR,
		}
		return
	}

	switch resp.StatusCode() {
	case http.StatusForbidden, http.StatusUnauthorized:
		cli.setIsRateLimit(true)
		errorChan <- &errorResp{
			Type: ERROR_MSG_AUTH_ERROR,
		}
	case http.StatusTooManyRequests:
		cli.setIsRateLimit(true)
		errorChan <- &errorResp{
			Type:   ERROR_MSG_RATE_LIMITED_ERROR,
			Status: resp.Status(),
		}
	case http.StatusOK:
		bodyChan <- resp.RawResponse
	case http.StatusServiceUnavailable, http.StatusInternalServerError, http.StatusBadGateway, http.StatusGatewayTimeout:
		errorChan <- &errorResp{
			Type:   ERROR_MSG_INTERNAL_ERROR,
			Status: resp.Status(),
		}
	default:
		errorChan <- &errorResp{
			Type:   ERROR_MSG_UNKNOWN_ERROR,
			Status: resp.Status(),
		}
	}
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
			eventChan <- event.Data()
		}
	}()

	var res ResponseBody
	for chunk := range eventChan {
		if err := json.Unmarshal([]byte(chunk), &res); err != nil {
			continue
		}
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
