package handler

import (
	"context"

	"notionboy/db/ent"
	"notionboy/internal/pkg/config"
)

func getAccFromContext(ctx context.Context) *ent.Account {
	acc := ctx.Value(config.ContextKeyUserAccount)
	// logger.SugaredLogger.Debugw("Get account from context", "acc", acc)
	if acc == nil {
		return nil
	}
	return acc.(*ent.Account)
}
