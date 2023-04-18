// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"notionboy/db/ent/predicate"
	"notionboy/db/ent/product"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// ProductUpdate is the builder for updating Product entities.
type ProductUpdate struct {
	config
	hooks    []Hook
	mutation *ProductMutation
}

// Where appends a list predicates to the ProductUpdate builder.
func (pu *ProductUpdate) Where(ps ...predicate.Product) *ProductUpdate {
	pu.mutation.Where(ps...)
	return pu
}

// SetUpdatedAt sets the "updated_at" field.
func (pu *ProductUpdate) SetUpdatedAt(t time.Time) *ProductUpdate {
	pu.mutation.SetUpdatedAt(t)
	return pu
}

// SetDeleted sets the "deleted" field.
func (pu *ProductUpdate) SetDeleted(b bool) *ProductUpdate {
	pu.mutation.SetDeleted(b)
	return pu
}

// SetNillableDeleted sets the "deleted" field if the given value is not nil.
func (pu *ProductUpdate) SetNillableDeleted(b *bool) *ProductUpdate {
	if b != nil {
		pu.SetDeleted(*b)
	}
	return pu
}

// SetName sets the "name" field.
func (pu *ProductUpdate) SetName(s string) *ProductUpdate {
	pu.mutation.SetName(s)
	return pu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (pu *ProductUpdate) SetNillableName(s *string) *ProductUpdate {
	if s != nil {
		pu.SetName(*s)
	}
	return pu
}

// SetDescription sets the "description" field.
func (pu *ProductUpdate) SetDescription(s string) *ProductUpdate {
	pu.mutation.SetDescription(s)
	return pu
}

// SetPrice sets the "price" field.
func (pu *ProductUpdate) SetPrice(f float64) *ProductUpdate {
	pu.mutation.ResetPrice()
	pu.mutation.SetPrice(f)
	return pu
}

// AddPrice adds f to the "price" field.
func (pu *ProductUpdate) AddPrice(f float64) *ProductUpdate {
	pu.mutation.AddPrice(f)
	return pu
}

// SetToken sets the "token" field.
func (pu *ProductUpdate) SetToken(i int64) *ProductUpdate {
	pu.mutation.ResetToken()
	pu.mutation.SetToken(i)
	return pu
}

// SetNillableToken sets the "token" field if the given value is not nil.
func (pu *ProductUpdate) SetNillableToken(i *int64) *ProductUpdate {
	if i != nil {
		pu.SetToken(*i)
	}
	return pu
}

// AddToken adds i to the "token" field.
func (pu *ProductUpdate) AddToken(i int64) *ProductUpdate {
	pu.mutation.AddToken(i)
	return pu
}

// SetStorage sets the "storage" field.
func (pu *ProductUpdate) SetStorage(i int64) *ProductUpdate {
	pu.mutation.ResetStorage()
	pu.mutation.SetStorage(i)
	return pu
}

// SetNillableStorage sets the "storage" field if the given value is not nil.
func (pu *ProductUpdate) SetNillableStorage(i *int64) *ProductUpdate {
	if i != nil {
		pu.SetStorage(*i)
	}
	return pu
}

// AddStorage adds i to the "storage" field.
func (pu *ProductUpdate) AddStorage(i int64) *ProductUpdate {
	pu.mutation.AddStorage(i)
	return pu
}

// Mutation returns the ProductMutation object of the builder.
func (pu *ProductUpdate) Mutation() *ProductMutation {
	return pu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (pu *ProductUpdate) Save(ctx context.Context) (int, error) {
	pu.defaults()
	return withHooks[int, ProductMutation](ctx, pu.sqlSave, pu.mutation, pu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (pu *ProductUpdate) SaveX(ctx context.Context) int {
	affected, err := pu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (pu *ProductUpdate) Exec(ctx context.Context) error {
	_, err := pu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pu *ProductUpdate) ExecX(ctx context.Context) {
	if err := pu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (pu *ProductUpdate) defaults() {
	if _, ok := pu.mutation.UpdatedAt(); !ok {
		v := product.UpdateDefaultUpdatedAt()
		pu.mutation.SetUpdatedAt(v)
	}
}

func (pu *ProductUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(product.Table, product.Columns, sqlgraph.NewFieldSpec(product.FieldID, field.TypeInt))
	if ps := pu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pu.mutation.UpdatedAt(); ok {
		_spec.SetField(product.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := pu.mutation.Deleted(); ok {
		_spec.SetField(product.FieldDeleted, field.TypeBool, value)
	}
	if value, ok := pu.mutation.Name(); ok {
		_spec.SetField(product.FieldName, field.TypeString, value)
	}
	if value, ok := pu.mutation.Description(); ok {
		_spec.SetField(product.FieldDescription, field.TypeString, value)
	}
	if value, ok := pu.mutation.Price(); ok {
		_spec.SetField(product.FieldPrice, field.TypeFloat64, value)
	}
	if value, ok := pu.mutation.AddedPrice(); ok {
		_spec.AddField(product.FieldPrice, field.TypeFloat64, value)
	}
	if value, ok := pu.mutation.Token(); ok {
		_spec.SetField(product.FieldToken, field.TypeInt64, value)
	}
	if value, ok := pu.mutation.AddedToken(); ok {
		_spec.AddField(product.FieldToken, field.TypeInt64, value)
	}
	if value, ok := pu.mutation.Storage(); ok {
		_spec.SetField(product.FieldStorage, field.TypeInt64, value)
	}
	if value, ok := pu.mutation.AddedStorage(); ok {
		_spec.AddField(product.FieldStorage, field.TypeInt64, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, pu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{product.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	pu.mutation.done = true
	return n, nil
}

// ProductUpdateOne is the builder for updating a single Product entity.
type ProductUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ProductMutation
}

// SetUpdatedAt sets the "updated_at" field.
func (puo *ProductUpdateOne) SetUpdatedAt(t time.Time) *ProductUpdateOne {
	puo.mutation.SetUpdatedAt(t)
	return puo
}

// SetDeleted sets the "deleted" field.
func (puo *ProductUpdateOne) SetDeleted(b bool) *ProductUpdateOne {
	puo.mutation.SetDeleted(b)
	return puo
}

// SetNillableDeleted sets the "deleted" field if the given value is not nil.
func (puo *ProductUpdateOne) SetNillableDeleted(b *bool) *ProductUpdateOne {
	if b != nil {
		puo.SetDeleted(*b)
	}
	return puo
}

// SetName sets the "name" field.
func (puo *ProductUpdateOne) SetName(s string) *ProductUpdateOne {
	puo.mutation.SetName(s)
	return puo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (puo *ProductUpdateOne) SetNillableName(s *string) *ProductUpdateOne {
	if s != nil {
		puo.SetName(*s)
	}
	return puo
}

// SetDescription sets the "description" field.
func (puo *ProductUpdateOne) SetDescription(s string) *ProductUpdateOne {
	puo.mutation.SetDescription(s)
	return puo
}

// SetPrice sets the "price" field.
func (puo *ProductUpdateOne) SetPrice(f float64) *ProductUpdateOne {
	puo.mutation.ResetPrice()
	puo.mutation.SetPrice(f)
	return puo
}

// AddPrice adds f to the "price" field.
func (puo *ProductUpdateOne) AddPrice(f float64) *ProductUpdateOne {
	puo.mutation.AddPrice(f)
	return puo
}

// SetToken sets the "token" field.
func (puo *ProductUpdateOne) SetToken(i int64) *ProductUpdateOne {
	puo.mutation.ResetToken()
	puo.mutation.SetToken(i)
	return puo
}

// SetNillableToken sets the "token" field if the given value is not nil.
func (puo *ProductUpdateOne) SetNillableToken(i *int64) *ProductUpdateOne {
	if i != nil {
		puo.SetToken(*i)
	}
	return puo
}

// AddToken adds i to the "token" field.
func (puo *ProductUpdateOne) AddToken(i int64) *ProductUpdateOne {
	puo.mutation.AddToken(i)
	return puo
}

// SetStorage sets the "storage" field.
func (puo *ProductUpdateOne) SetStorage(i int64) *ProductUpdateOne {
	puo.mutation.ResetStorage()
	puo.mutation.SetStorage(i)
	return puo
}

// SetNillableStorage sets the "storage" field if the given value is not nil.
func (puo *ProductUpdateOne) SetNillableStorage(i *int64) *ProductUpdateOne {
	if i != nil {
		puo.SetStorage(*i)
	}
	return puo
}

// AddStorage adds i to the "storage" field.
func (puo *ProductUpdateOne) AddStorage(i int64) *ProductUpdateOne {
	puo.mutation.AddStorage(i)
	return puo
}

// Mutation returns the ProductMutation object of the builder.
func (puo *ProductUpdateOne) Mutation() *ProductMutation {
	return puo.mutation
}

// Where appends a list predicates to the ProductUpdate builder.
func (puo *ProductUpdateOne) Where(ps ...predicate.Product) *ProductUpdateOne {
	puo.mutation.Where(ps...)
	return puo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (puo *ProductUpdateOne) Select(field string, fields ...string) *ProductUpdateOne {
	puo.fields = append([]string{field}, fields...)
	return puo
}

// Save executes the query and returns the updated Product entity.
func (puo *ProductUpdateOne) Save(ctx context.Context) (*Product, error) {
	puo.defaults()
	return withHooks[*Product, ProductMutation](ctx, puo.sqlSave, puo.mutation, puo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (puo *ProductUpdateOne) SaveX(ctx context.Context) *Product {
	node, err := puo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (puo *ProductUpdateOne) Exec(ctx context.Context) error {
	_, err := puo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (puo *ProductUpdateOne) ExecX(ctx context.Context) {
	if err := puo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (puo *ProductUpdateOne) defaults() {
	if _, ok := puo.mutation.UpdatedAt(); !ok {
		v := product.UpdateDefaultUpdatedAt()
		puo.mutation.SetUpdatedAt(v)
	}
}

func (puo *ProductUpdateOne) sqlSave(ctx context.Context) (_node *Product, err error) {
	_spec := sqlgraph.NewUpdateSpec(product.Table, product.Columns, sqlgraph.NewFieldSpec(product.FieldID, field.TypeInt))
	id, ok := puo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Product.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := puo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, product.FieldID)
		for _, f := range fields {
			if !product.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != product.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := puo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := puo.mutation.UpdatedAt(); ok {
		_spec.SetField(product.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := puo.mutation.Deleted(); ok {
		_spec.SetField(product.FieldDeleted, field.TypeBool, value)
	}
	if value, ok := puo.mutation.Name(); ok {
		_spec.SetField(product.FieldName, field.TypeString, value)
	}
	if value, ok := puo.mutation.Description(); ok {
		_spec.SetField(product.FieldDescription, field.TypeString, value)
	}
	if value, ok := puo.mutation.Price(); ok {
		_spec.SetField(product.FieldPrice, field.TypeFloat64, value)
	}
	if value, ok := puo.mutation.AddedPrice(); ok {
		_spec.AddField(product.FieldPrice, field.TypeFloat64, value)
	}
	if value, ok := puo.mutation.Token(); ok {
		_spec.SetField(product.FieldToken, field.TypeInt64, value)
	}
	if value, ok := puo.mutation.AddedToken(); ok {
		_spec.AddField(product.FieldToken, field.TypeInt64, value)
	}
	if value, ok := puo.mutation.Storage(); ok {
		_spec.SetField(product.FieldStorage, field.TypeInt64, value)
	}
	if value, ok := puo.mutation.AddedStorage(); ok {
		_spec.AddField(product.FieldStorage, field.TypeInt64, value)
	}
	_node = &Product{config: puo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, puo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{product.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	puo.mutation.done = true
	return _node, nil
}
