package dao

import (
	"context"

	"notionboy/db/ent"
	"notionboy/db/ent/prompt"
	"notionboy/internal/pkg/db"

	"github.com/google/uuid"
)

func ListPrompts(ctx context.Context, userId uuid.UUID) ([]*ent.Prompt, error) {
	return db.GetClient().Prompt.Query().
		Where(prompt.UserIDEQ(userId)).
		All(ctx)
}

func GetPrompt(ctx context.Context, id, userId uuid.UUID) (*ent.Prompt, error) {
	return db.GetClient().Prompt.Query().
		Where(prompt.UUIDEQ(id), prompt.UserIDEQ(userId)).
		Only(ctx)
}

func CreatePrompt(ctx context.Context, userId uuid.UUID, act, prompt string) (*ent.Prompt, error) {
	return db.GetClient().Prompt.Create().
		SetUUID(uuid.New()).
		SetUserID(userId).
		SetAct(act).
		SetPrompt(prompt).
		SetIsCustom(true).
		Save(ctx)
}

func UpdatePrompt(ctx context.Context, id, userId uuid.UUID, act, promptStr string) (*ent.Prompt, error) {
	if _, err := db.GetClient().Prompt.Update().
		SetAct(act).
		SetPrompt(promptStr).
		Where(prompt.UUIDEQ(id), prompt.UserIDEQ(userId)).
		Save(ctx); err != nil {
		return nil, err
	} else {
		return GetPrompt(ctx, id, userId)
	}
}

func DeletePrompt(ctx context.Context, id, userId uuid.UUID) error {
	_, err := db.GetClient().Prompt.Delete().
		Where(prompt.UUIDEQ(id), prompt.UserIDEQ(userId)).
		Exec(ctx)
	return err
}
