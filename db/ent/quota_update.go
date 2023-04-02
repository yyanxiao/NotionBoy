// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"notionboy/db/ent/predicate"
	"notionboy/db/ent/quota"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// QuotaUpdate is the builder for updating Quota entities.
type QuotaUpdate struct {
	config
	hooks    []Hook
	mutation *QuotaMutation
}

// Where appends a list predicates to the QuotaUpdate builder.
func (qu *QuotaUpdate) Where(ps ...predicate.Quota) *QuotaUpdate {
	qu.mutation.Where(ps...)
	return qu
}

// SetUpdatedAt sets the "updated_at" field.
func (qu *QuotaUpdate) SetUpdatedAt(t time.Time) *QuotaUpdate {
	qu.mutation.SetUpdatedAt(t)
	return qu
}

// SetDeleted sets the "deleted" field.
func (qu *QuotaUpdate) SetDeleted(b bool) *QuotaUpdate {
	qu.mutation.SetDeleted(b)
	return qu
}

// SetNillableDeleted sets the "deleted" field if the given value is not nil.
func (qu *QuotaUpdate) SetNillableDeleted(b *bool) *QuotaUpdate {
	if b != nil {
		qu.SetDeleted(*b)
	}
	return qu
}

// SetUserID sets the "user_id" field.
func (qu *QuotaUpdate) SetUserID(i int) *QuotaUpdate {
	qu.mutation.ResetUserID()
	qu.mutation.SetUserID(i)
	return qu
}

// AddUserID adds i to the "user_id" field.
func (qu *QuotaUpdate) AddUserID(i int) *QuotaUpdate {
	qu.mutation.AddUserID(i)
	return qu
}

// SetPlan sets the "plan" field.
func (qu *QuotaUpdate) SetPlan(s string) *QuotaUpdate {
	qu.mutation.SetPlan(s)
	return qu
}

// SetResetTime sets the "reset_time" field.
func (qu *QuotaUpdate) SetResetTime(t time.Time) *QuotaUpdate {
	qu.mutation.SetResetTime(t)
	return qu
}

// SetToken sets the "token" field.
func (qu *QuotaUpdate) SetToken(i int64) *QuotaUpdate {
	qu.mutation.ResetToken()
	qu.mutation.SetToken(i)
	return qu
}

// AddToken adds i to the "token" field.
func (qu *QuotaUpdate) AddToken(i int64) *QuotaUpdate {
	qu.mutation.AddToken(i)
	return qu
}

// SetTokenUsed sets the "token_used" field.
func (qu *QuotaUpdate) SetTokenUsed(i int64) *QuotaUpdate {
	qu.mutation.ResetTokenUsed()
	qu.mutation.SetTokenUsed(i)
	return qu
}

// SetNillableTokenUsed sets the "token_used" field if the given value is not nil.
func (qu *QuotaUpdate) SetNillableTokenUsed(i *int64) *QuotaUpdate {
	if i != nil {
		qu.SetTokenUsed(*i)
	}
	return qu
}

// AddTokenUsed adds i to the "token_used" field.
func (qu *QuotaUpdate) AddTokenUsed(i int64) *QuotaUpdate {
	qu.mutation.AddTokenUsed(i)
	return qu
}

// Mutation returns the QuotaMutation object of the builder.
func (qu *QuotaUpdate) Mutation() *QuotaMutation {
	return qu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (qu *QuotaUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	qu.defaults()
	if len(qu.hooks) == 0 {
		affected, err = qu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*QuotaMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			qu.mutation = mutation
			affected, err = qu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(qu.hooks) - 1; i >= 0; i-- {
			if qu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = qu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, qu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (qu *QuotaUpdate) SaveX(ctx context.Context) int {
	affected, err := qu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (qu *QuotaUpdate) Exec(ctx context.Context) error {
	_, err := qu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (qu *QuotaUpdate) ExecX(ctx context.Context) {
	if err := qu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (qu *QuotaUpdate) defaults() {
	if _, ok := qu.mutation.UpdatedAt(); !ok {
		v := quota.UpdateDefaultUpdatedAt()
		qu.mutation.SetUpdatedAt(v)
	}
}

func (qu *QuotaUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   quota.Table,
			Columns: quota.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: quota.FieldID,
			},
		},
	}
	if ps := qu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := qu.mutation.UpdatedAt(); ok {
		_spec.SetField(quota.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := qu.mutation.Deleted(); ok {
		_spec.SetField(quota.FieldDeleted, field.TypeBool, value)
	}
	if value, ok := qu.mutation.UserID(); ok {
		_spec.SetField(quota.FieldUserID, field.TypeInt, value)
	}
	if value, ok := qu.mutation.AddedUserID(); ok {
		_spec.AddField(quota.FieldUserID, field.TypeInt, value)
	}
	if value, ok := qu.mutation.Plan(); ok {
		_spec.SetField(quota.FieldPlan, field.TypeString, value)
	}
	if value, ok := qu.mutation.ResetTime(); ok {
		_spec.SetField(quota.FieldResetTime, field.TypeTime, value)
	}
	if value, ok := qu.mutation.Token(); ok {
		_spec.SetField(quota.FieldToken, field.TypeInt64, value)
	}
	if value, ok := qu.mutation.AddedToken(); ok {
		_spec.AddField(quota.FieldToken, field.TypeInt64, value)
	}
	if value, ok := qu.mutation.TokenUsed(); ok {
		_spec.SetField(quota.FieldTokenUsed, field.TypeInt64, value)
	}
	if value, ok := qu.mutation.AddedTokenUsed(); ok {
		_spec.AddField(quota.FieldTokenUsed, field.TypeInt64, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, qu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{quota.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	return n, nil
}

// QuotaUpdateOne is the builder for updating a single Quota entity.
type QuotaUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *QuotaMutation
}

// SetUpdatedAt sets the "updated_at" field.
func (quo *QuotaUpdateOne) SetUpdatedAt(t time.Time) *QuotaUpdateOne {
	quo.mutation.SetUpdatedAt(t)
	return quo
}

// SetDeleted sets the "deleted" field.
func (quo *QuotaUpdateOne) SetDeleted(b bool) *QuotaUpdateOne {
	quo.mutation.SetDeleted(b)
	return quo
}

// SetNillableDeleted sets the "deleted" field if the given value is not nil.
func (quo *QuotaUpdateOne) SetNillableDeleted(b *bool) *QuotaUpdateOne {
	if b != nil {
		quo.SetDeleted(*b)
	}
	return quo
}

// SetUserID sets the "user_id" field.
func (quo *QuotaUpdateOne) SetUserID(i int) *QuotaUpdateOne {
	quo.mutation.ResetUserID()
	quo.mutation.SetUserID(i)
	return quo
}

// AddUserID adds i to the "user_id" field.
func (quo *QuotaUpdateOne) AddUserID(i int) *QuotaUpdateOne {
	quo.mutation.AddUserID(i)
	return quo
}

// SetPlan sets the "plan" field.
func (quo *QuotaUpdateOne) SetPlan(s string) *QuotaUpdateOne {
	quo.mutation.SetPlan(s)
	return quo
}

// SetResetTime sets the "reset_time" field.
func (quo *QuotaUpdateOne) SetResetTime(t time.Time) *QuotaUpdateOne {
	quo.mutation.SetResetTime(t)
	return quo
}

// SetToken sets the "token" field.
func (quo *QuotaUpdateOne) SetToken(i int64) *QuotaUpdateOne {
	quo.mutation.ResetToken()
	quo.mutation.SetToken(i)
	return quo
}

// AddToken adds i to the "token" field.
func (quo *QuotaUpdateOne) AddToken(i int64) *QuotaUpdateOne {
	quo.mutation.AddToken(i)
	return quo
}

// SetTokenUsed sets the "token_used" field.
func (quo *QuotaUpdateOne) SetTokenUsed(i int64) *QuotaUpdateOne {
	quo.mutation.ResetTokenUsed()
	quo.mutation.SetTokenUsed(i)
	return quo
}

// SetNillableTokenUsed sets the "token_used" field if the given value is not nil.
func (quo *QuotaUpdateOne) SetNillableTokenUsed(i *int64) *QuotaUpdateOne {
	if i != nil {
		quo.SetTokenUsed(*i)
	}
	return quo
}

// AddTokenUsed adds i to the "token_used" field.
func (quo *QuotaUpdateOne) AddTokenUsed(i int64) *QuotaUpdateOne {
	quo.mutation.AddTokenUsed(i)
	return quo
}

// Mutation returns the QuotaMutation object of the builder.
func (quo *QuotaUpdateOne) Mutation() *QuotaMutation {
	return quo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (quo *QuotaUpdateOne) Select(field string, fields ...string) *QuotaUpdateOne {
	quo.fields = append([]string{field}, fields...)
	return quo
}

// Save executes the query and returns the updated Quota entity.
func (quo *QuotaUpdateOne) Save(ctx context.Context) (*Quota, error) {
	var (
		err  error
		node *Quota
	)
	quo.defaults()
	if len(quo.hooks) == 0 {
		node, err = quo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*QuotaMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			quo.mutation = mutation
			node, err = quo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(quo.hooks) - 1; i >= 0; i-- {
			if quo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = quo.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, quo.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*Quota)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from QuotaMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (quo *QuotaUpdateOne) SaveX(ctx context.Context) *Quota {
	node, err := quo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (quo *QuotaUpdateOne) Exec(ctx context.Context) error {
	_, err := quo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (quo *QuotaUpdateOne) ExecX(ctx context.Context) {
	if err := quo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (quo *QuotaUpdateOne) defaults() {
	if _, ok := quo.mutation.UpdatedAt(); !ok {
		v := quota.UpdateDefaultUpdatedAt()
		quo.mutation.SetUpdatedAt(v)
	}
}

func (quo *QuotaUpdateOne) sqlSave(ctx context.Context) (_node *Quota, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   quota.Table,
			Columns: quota.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: quota.FieldID,
			},
		},
	}
	id, ok := quo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Quota.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := quo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, quota.FieldID)
		for _, f := range fields {
			if !quota.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != quota.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := quo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := quo.mutation.UpdatedAt(); ok {
		_spec.SetField(quota.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := quo.mutation.Deleted(); ok {
		_spec.SetField(quota.FieldDeleted, field.TypeBool, value)
	}
	if value, ok := quo.mutation.UserID(); ok {
		_spec.SetField(quota.FieldUserID, field.TypeInt, value)
	}
	if value, ok := quo.mutation.AddedUserID(); ok {
		_spec.AddField(quota.FieldUserID, field.TypeInt, value)
	}
	if value, ok := quo.mutation.Plan(); ok {
		_spec.SetField(quota.FieldPlan, field.TypeString, value)
	}
	if value, ok := quo.mutation.ResetTime(); ok {
		_spec.SetField(quota.FieldResetTime, field.TypeTime, value)
	}
	if value, ok := quo.mutation.Token(); ok {
		_spec.SetField(quota.FieldToken, field.TypeInt64, value)
	}
	if value, ok := quo.mutation.AddedToken(); ok {
		_spec.AddField(quota.FieldToken, field.TypeInt64, value)
	}
	if value, ok := quo.mutation.TokenUsed(); ok {
		_spec.SetField(quota.FieldTokenUsed, field.TypeInt64, value)
	}
	if value, ok := quo.mutation.AddedTokenUsed(); ok {
		_spec.AddField(quota.FieldTokenUsed, field.TypeInt64, value)
	}
	_node = &Quota{config: quo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, quo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{quota.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	return _node, nil
}
