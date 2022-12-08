package chatgpt

const chatURL = "https://chat.openai.com/backend-api/conversation"

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
