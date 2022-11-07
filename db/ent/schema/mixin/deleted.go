package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type DeletedMixin struct {
	mixin.Schema
}

func (DeletedMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("deleted").
			Default(false),
	}
}
