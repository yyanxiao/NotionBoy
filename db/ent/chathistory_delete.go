// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"notionboy/db/ent/chathistory"
	"notionboy/db/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// ChatHistoryDelete is the builder for deleting a ChatHistory entity.
type ChatHistoryDelete struct {
	config
	hooks    []Hook
	mutation *ChatHistoryMutation
}

// Where appends a list predicates to the ChatHistoryDelete builder.
func (chd *ChatHistoryDelete) Where(ps ...predicate.ChatHistory) *ChatHistoryDelete {
	chd.mutation.Where(ps...)
	return chd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (chd *ChatHistoryDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(chd.hooks) == 0 {
		affected, err = chd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ChatHistoryMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			chd.mutation = mutation
			affected, err = chd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(chd.hooks) - 1; i >= 0; i-- {
			if chd.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = chd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, chd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (chd *ChatHistoryDelete) ExecX(ctx context.Context) int {
	n, err := chd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (chd *ChatHistoryDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: chathistory.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: chathistory.FieldID,
			},
		},
	}
	if ps := chd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, chd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	return affected, err
}

// ChatHistoryDeleteOne is the builder for deleting a single ChatHistory entity.
type ChatHistoryDeleteOne struct {
	chd *ChatHistoryDelete
}

// Exec executes the deletion query.
func (chdo *ChatHistoryDeleteOne) Exec(ctx context.Context) error {
	n, err := chdo.chd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{chathistory.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (chdo *ChatHistoryDeleteOne) ExecX(ctx context.Context) {
	chdo.chd.ExecX(ctx)
}