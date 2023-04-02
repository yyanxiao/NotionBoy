package dao

import (
	"context"
	"notionboy/db/ent"
	"notionboy/db/ent/order"
	"notionboy/db/ent/product"
	"notionboy/internal/pkg/db"

	"github.com/google/uuid"
)

func CreateProduct(ctx context.Context, name string, price float64, token int64) (*ent.Product, error) {
	return db.GetClient().Product.Create().
		SetUUID(uuid.New()).
		SetName(name).
		SetToken(token).
		SetPrice(price).
		Save(ctx)
}

func ListProductsByIds(ctx context.Context, ids []uuid.UUID) ([]*ent.Product, error) {
	query := db.GetClient().Product.Query()
	return query.
		Where(product.UUIDIn(ids...)).
		All(ctx)
}

func ListProducts(ctx context.Context, limit, offset int) ([]*ent.Product, error) {
	return db.GetClient().Product.Query().
		Order(ent.Desc(order.FieldCreatedAt)).
		Limit(limit).
		Offset(offset).
		All(ctx)
}

func GetProduct(ctx context.Context, id uuid.UUID) (*ent.Product, error) {
	return db.GetClient().Product.Query().
		Where(product.UUIDEQ(id)).
		Only(ctx)
}

func DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := db.GetClient().Product.Delete().
		Where(product.UUIDEQ(id)).
		Exec(ctx)
	return err
}

func UpdateProduct(ctx context.Context, id uuid.UUID, name string, price float64, token int64) (*ent.Product, error) {
	query := db.GetClient().Product.Update().
		Where(product.UUIDEQ(id))

	if name != "" {
		query = query.SetName(name)
	}

	if price != 0 {
		query = query.SetPrice(price)
	}

	if token != 0 {
		query = query.SetToken(token)
	}

	if err := query.Exec(ctx); err != nil {
		return nil, err
	}

	return GetProduct(ctx, id)
}
