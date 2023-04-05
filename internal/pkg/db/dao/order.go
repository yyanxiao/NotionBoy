package dao

import (
	"context"
	"strings"

	"notionboy/db/ent"
	"notionboy/db/ent/order"
	"notionboy/internal/pkg/db"

	"github.com/google/uuid"
)

func CreateOrder(ctx context.Context, userId, uid uuid.UUID, price float64, note string, tokens int64) (*ent.Order, error) {
	if uid == uuid.Nil {
		uid = uuid.New()
	}
	return db.GetClient().Order.Create().
		SetUUID(uid).
		SetUserID(userId).
		SetPrice(price).
		SetNote(strings.TrimSpace(note)).
		Save(ctx)
}

func ListOrders(ctx context.Context, userId uuid.UUID, status order.Status, limit, offset int) ([]*ent.Order, error) {
	query := db.GetClient().Order.Query().
		Where(order.UserIDEQ(userId))

	if status != "" {
		query = query.Where(order.StatusEQ(status))
	}

	return query.
		Order(ent.Desc(order.FieldCreatedAt)).
		Limit(limit).
		Offset(offset).
		All(ctx)
}

func GetOrder(ctx context.Context, id, userId uuid.UUID) (*ent.Order, error) {
	return db.GetClient().Order.Query().
		Where(order.UUIDEQ(id)).
		Where(order.UserIDEQ(userId)).
		Only(ctx)
}

func DeleteOrder(ctx context.Context, id, userId uuid.UUID) error {
	_, err := db.GetClient().Order.Delete().
		Where(order.UUIDEQ(id)).
		Where(order.UserIDEQ(userId)).
		Exec(ctx)
	return err
}

func UpdateOrder(ctx context.Context, id, userId uuid.UUID, status order.Status, price float64, note, paymentInfo string, tokens int64) (*ent.Order, error) {
	query := db.GetClient().Order.Update().
		Where(order.UUIDEQ(id), order.UserIDEQ(userId)).
		SetStatus(status)

	if status != "" {
		query = query.SetStatus(status)
	}
	if price > 0 {
		query = query.SetPrice(price)
	}
	if note != "" {
		query = query.SetNote(note)
	}
	if paymentInfo != "" {
		query = query.SetPaymentInfo(paymentInfo)
	}

	if err := query.Exec(ctx); err != nil {
		return nil, err
	}

	return GetOrder(ctx, id, userId)
}

func UpdateOrderStatus(cli *ent.Client, ctx context.Context, id, userId uuid.UUID, status order.Status) error {
	return cli.Order.Update().
		Where(order.UUIDEQ(id), order.UserIDEQ(userId)).
		SetStatus(status).
		Exec(ctx)
}
