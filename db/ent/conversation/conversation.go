// Code generated by ent, DO NOT EDIT.

package conversation

import (
	"time"
)

const (
	// Label holds the string label denoting the conversation type in the database.
	Label = "conversation"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldDeleted holds the string denoting the deleted field in the database.
	FieldDeleted = "deleted"
	// FieldUUID holds the string denoting the uuid field in the database.
	FieldUUID = "uuid"
	// FieldUserID holds the string denoting the user_id field in the database.
	FieldUserID = "user_id"
	// FieldInstruction holds the string denoting the instruction field in the database.
	FieldInstruction = "instruction"
	// FieldTitle holds the string denoting the title field in the database.
	FieldTitle = "title"
	// FieldTokenUsage holds the string denoting the token_usage field in the database.
	FieldTokenUsage = "token_usage"
	// EdgeConversationMessages holds the string denoting the conversation_messages edge name in mutations.
	EdgeConversationMessages = "conversation_messages"
	// Table holds the table name of the conversation in the database.
	Table = "conversations"
	// ConversationMessagesTable is the table that holds the conversation_messages relation/edge.
	ConversationMessagesTable = "conversation_messages"
	// ConversationMessagesInverseTable is the table name for the ConversationMessage entity.
	// It exists in this package in order to avoid circular dependency with the "conversationmessage" package.
	ConversationMessagesInverseTable = "conversation_messages"
	// ConversationMessagesColumn is the table column denoting the conversation_messages relation/edge.
	ConversationMessagesColumn = "conversation_conversation_messages"
)

// Columns holds all SQL columns for conversation fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldDeleted,
	FieldUUID,
	FieldUserID,
	FieldInstruction,
	FieldTitle,
	FieldTokenUsage,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
	// DefaultDeleted holds the default value on creation for the "deleted" field.
	DefaultDeleted bool
	// DefaultTokenUsage holds the default value on creation for the "token_usage" field.
	DefaultTokenUsage int64
)
