package dao

import (
	"context"
	"notionboy/db/ent"
	"notionboy/internal/pkg/db"
)

func GenTx(ctx context.Context) (*ent.Tx, error) {
	return db.GetClient().Tx(ctx)
}

// func getTx(ctx context.Context) (*ent.Tx, error) {
// 	val := ctx.Value(config.ContentKeyTransaction)
// 	if val == nil {
// 		return db.GetClient().Tx(ctx)
// 	}

// 	return val.(*ent.Tx), nil
// }
