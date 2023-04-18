package schema

import (
	"notionboy/db/ent/schema/mixin"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type Order struct {
	ent.Schema
}

func (Order) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("uuid", uuid.UUID{}).Unique().Immutable().Comment("UUID"),
		field.UUID("user_id", uuid.UUID{}).Comment("user id"),
		field.UUID("product_id", uuid.UUID{}).Comment("product id"),
		field.Float("price").Comment("total price"),
		field.Enum("status").Values("Unpaid", "Paying", "Paid", "Processing", "Cancelled", "Refunded", "Completed").Default("Unpaid"),
		field.Text("note").Optional().Comment("Note for the order"),
		field.Text("payment_info").Optional().Comment("Payment info"),
	}
}

func (Order) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
		mixin.DeletedMixin{},
	}
}

func (Order) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "created_at"),
	}
}
