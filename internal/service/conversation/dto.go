package conversation

import (
	"time"

	"notionboy/api/pb/model"
	"notionboy/db/ent"

	"github.com/google/uuid"
)

type ConversationDTO struct {
	ID          string                    `json:"id"`
	CreatedAt   time.Time                 `json:"created_at"`
	UpdatedAt   time.Time                 `json:"updated_at"`
	Deleted     bool                      `json:"deleted"`
	UserID      string                    `json:"user_id"`
	Instruction string                    `json:"instruction"`
	Title       string                    `json:"title"`
	Messages    []*ConversationMessageDTO `json:"messages"`
}

func (d *ConversationDTO) ToDB() *ent.Conversation {
	return &ent.Conversation{
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
		Deleted:     d.Deleted,
		UUID:        uuid.MustParse(d.ID),
		UserID:      uuid.MustParse(d.UserID),
		Instruction: d.Instruction,
		Title:       d.Title,
	}
}

func (d *ConversationDTO) ToPB() *model.Conversation {
	messages := make([]*model.Message, 0, len(d.Messages))
	for _, message := range d.Messages {
		messages = append(messages, message.ToPB())
	}
	return &model.Conversation{
		Id:          d.ID,
		CreatedAt:   d.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   d.UpdatedAt.Format(time.RFC3339),
		Instruction: d.Instruction,
		Title:       d.Title,
		Messages:    messages,
	}
}

func (d *ConversationDTO) ToPBWithoutMessages() *model.ConversationWithoutMessages {
	return &model.ConversationWithoutMessages{
		Id:          d.ID,
		CreatedAt:   d.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   d.UpdatedAt.Format(time.RFC3339),
		Instruction: d.Instruction,
		Title:       d.Title,
	}
}

func (d *ConversationDTO) FromDB(c *ent.Conversation) *ConversationDTO {
	// logger.SugaredLogger.Debugw("ConversationDTO.FromDB", "conversation", c)
	d.CreatedAt = c.CreatedAt
	d.UpdatedAt = c.UpdatedAt
	d.Deleted = c.Deleted
	d.ID = c.UUID.String()
	d.UserID = c.UserID.String()
	d.Instruction = c.Instruction
	d.Title = c.Title
	return d
}

func ConversationDTOFromDB(c *ent.Conversation) *ConversationDTO {
	return (&ConversationDTO{}).FromDB(c)
}

type ConversationMessageDTO struct {
	ID             string    `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Deleted        bool      `json:"deleted"`
	UserID         string    `json:"user_id"`
	ConversationID string    `json:"conversation_id"`
	Request        string    `json:"request"`
	Response       string    `json:"response"`
	TokenUsage     int64     `json:"token_usage"`
}

func (d *ConversationMessageDTO) ToDB() *ent.ConversationMessage {
	return &ent.ConversationMessage{
		CreatedAt:      d.CreatedAt,
		UpdatedAt:      d.UpdatedAt,
		Deleted:        d.Deleted,
		UUID:           uuid.MustParse(d.ID),
		UserID:         uuid.MustParse(d.UserID),
		ConversationID: uuid.MustParse(d.ConversationID),
		Request:        d.Request,
		Response:       d.Response,
		TokenUsage:     d.TokenUsage,
	}
}

func (d *ConversationMessageDTO) FromDB(m *ent.ConversationMessage) *ConversationMessageDTO {
	d.CreatedAt = m.CreatedAt
	d.UpdatedAt = m.UpdatedAt
	d.Deleted = m.Deleted
	d.ID = m.UUID.String()
	d.UserID = m.UserID.String()
	d.ConversationID = m.ConversationID.String()
	d.Request = m.Request
	d.Response = m.Response
	d.TokenUsage = m.TokenUsage
	return d
}

func (d *ConversationMessageDTO) ToPB() *model.Message {
	return &model.Message{
		Id:             d.ID,
		ConversationId: d.ConversationID,
		CreatedAt:      d.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      d.UpdatedAt.Format(time.RFC3339),
		Request:        d.Request,
		Response:       d.Response,
		TokenUsage:     int32(d.TokenUsage),
	}
}

func ConversationMessageDTOFromDB(m *ent.ConversationMessage) *ConversationMessageDTO {
	return (&ConversationMessageDTO{}).FromDB(m)
}
