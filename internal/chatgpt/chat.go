package chatgpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"notionboy/internal/pkg/logger"
	"sync"

	"github.com/google/uuid"
)

const chatURL = "https://chat.openai.com/backend-api/conversation"

var mu sync.Mutex

type MessageContent struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

type Message struct {
	ID      string         `json:"id"`
	Role    string         `json:"role"`
	Content MessageContent `json:"content"`
}

type Payload struct {
	Action          string    `json:"action"`
	Messages        []Message `json:"messages"`
	ConversationId  string    `json:"conversation_id,omitempty"`
	ParentMessageId string    `json:"parent_message_id"`
	Model           string    `json:"model"`
}

type ResponseBody struct {
	Message        Message `json:"message"`
	ConversationId string  `json:"conversation_id,omitempty"`
	Error          string  `json:"error"`
}

// Chat
// input: parentMessageId, conversationId , prompt
// output: parentMessageId , message, error
func Chat(parentMessageId, prompt string) (string, string, error) {
	// only enable one chat per time
	mu.Lock()
	defer mu.Unlock()
	payload := buildPayload(parentMessageId, []string{prompt})
	logger.SugaredLogger.Debugw("Request Payload", "payload", payload)

	resp, err := client.R().
		SetBody(payload).
		SetDoNotParseResponse(true).
		Post(chatURL)
	if err != nil {
		logger.SugaredLogger.Errorw("Talk to chatGPT error", "err", err)
		return "", "", err
	}

	logger.SugaredLogger.Debugw("status", "status_code", resp.Status())
	if resp.StatusCode() != http.StatusOK {
		logger.SugaredLogger.Errorw("Chat with chatGPT error, please retry first", "status_code", resp.StatusCode(), "resp_text", string(resp.Body()))
		RefreshSession()
		return "", "", errors.New("Chat with chatGPT error. \nerror: " + resp.Status())
	}
	newParentMessageId, _, messages, err := processResponse(resp.RawResponse)
	return newParentMessageId, messages[0], err
}

// processResponse
// resp: *resty.Response as input
// parentMessageId, conversationId, messages, error
func processResponse(resp *http.Response) (string, string, []string, error) {
	defer resp.Body.Close()
	//respBody, _ := io.ReadAll(resp.Body)
	//dataSlice := bytes.Split(respBody, []byte{'\n'})
	//var data []byte
	//for i := len(dataSlice) - 1; i >= 0; i-- {
	//	if strings.HasPrefix(string(dataSlice[i]), "data: {") {
	//		data = dataSlice[i]
	//		data = data[6:]
	//		break
	//	}
	//}
	//if data == nil {
	//	return "", "", []string{""}, errors.New("read response error, please retry")
	//}
	data := readBody(resp.Body)
	if data == nil {
		return "", "", []string{""}, errors.New("read response error, please retry")
	}

	var res ResponseBody
	_ = json.Unmarshal(data, &res)
	logger.SugaredLogger.Debugw("Response", "conversation_id", res.ConversationId, "error", res.Error, "message", res.Message)
	parentMessageId, conversationId, messages := res.Message.ID, res.ConversationId, res.Message.Content.Parts
	return parentMessageId, conversationId, messages, nil
}

func readBody(r io.Reader) []byte {
	b := make([]byte, 2)
	pending := make([]byte, 0)
	// lines := make([][]byte, 0)
	res := make([]byte, 0)
	dataPrefix := []byte("data: {")
	for {
		if _, err := r.Read(b); err != nil {
			if err == io.EOF {
				tmpLine := make([]byte, 0)
				for _, item := range pending {
					if item == 0 {
						break
					}
					tmpLine = append(tmpLine, item)
				}
				if bytes.HasPrefix(tmpLine, dataPrefix) {
					res = tmpLine
				}
				// lines = append(lines, tmpLine)
				break
			}
		}
		if len(pending) != 0 {
			b = append(pending[:], b...)
			pending = make([]byte, 0)
		}
		if len(b) > 0 {
			tmpLines := bytes.Split(b, []byte{'\n'})
			lastLine := tmpLines[len(tmpLines)-1]
			if len(lastLine) > 0 && lastLine[len(lastLine)-1] == b[len(b)-1] {
				pending = append(pending, lastLine...)
				tmpLines = tmpLines[:len(tmpLines)-1]
			} else {
				pending = make([]byte, 0)
			}
			for _, line := range tmpLines {
				if len(line) > 0 {
					if bytes.HasPrefix(line, dataPrefix) {
						res = line
					}
					// lines = append(lines, line)
				}
			}
		}
		b = make([]byte, 2)
	}
	if len(res) > 0 {
		return res[6:]
	}
	return res
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
