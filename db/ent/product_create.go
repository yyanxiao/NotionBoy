// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"notionboy/db/ent/product"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// ProductCreate is the builder for creating a Product entity.
type ProductCreate struct {
	config
	mutation *ProductMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreatedAt sets the "created_at" field.
func (pc *ProductCreate) SetCreatedAt(t time.Time) *ProductCreate {
	pc.mutation.SetCreatedAt(t)
	return pc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (pc *ProductCreate) SetNillableCreatedAt(t *time.Time) *ProductCreate {
	if t != nil {
		pc.SetCreatedAt(*t)
	}
	return pc
}

// SetUpdatedAt sets the "updated_at" field.
func (pc *ProductCreate) SetUpdatedAt(t time.Time) *ProductCreate {
	pc.mutation.SetUpdatedAt(t)
	return pc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (pc *ProductCreate) SetNillableUpdatedAt(t *time.Time) *ProductCreate {
	if t != nil {
		pc.SetUpdatedAt(*t)
	}
	return pc
}

// SetDeleted sets the "deleted" field.
func (pc *ProductCreate) SetDeleted(b bool) *ProductCreate {
	pc.mutation.SetDeleted(b)
	return pc
}

// SetNillableDeleted sets the "deleted" field if the given value is not nil.
func (pc *ProductCreate) SetNillableDeleted(b *bool) *ProductCreate {
	if b != nil {
		pc.SetDeleted(*b)
	}
	return pc
}

// SetUUID sets the "uuid" field.
func (pc *ProductCreate) SetUUID(u uuid.UUID) *ProductCreate {
	pc.mutation.SetUUID(u)
	return pc
}

// SetName sets the "name" field.
func (pc *ProductCreate) SetName(s string) *ProductCreate {
	pc.mutation.SetName(s)
	return pc
}

// SetNillableName sets the "name" field if the given value is not nil.
func (pc *ProductCreate) SetNillableName(s *string) *ProductCreate {
	if s != nil {
		pc.SetName(*s)
	}
	return pc
}

// SetDescription sets the "description" field.
func (pc *ProductCreate) SetDescription(s string) *ProductCreate {
	pc.mutation.SetDescription(s)
	return pc
}

// SetPrice sets the "price" field.
func (pc *ProductCreate) SetPrice(f float64) *ProductCreate {
	pc.mutation.SetPrice(f)
	return pc
}

// SetToken sets the "token" field.
func (pc *ProductCreate) SetToken(i int64) *ProductCreate {
	pc.mutation.SetToken(i)
	return pc
}

// SetNillableToken sets the "token" field if the given value is not nil.
func (pc *ProductCreate) SetNillableToken(i *int64) *ProductCreate {
	if i != nil {
		pc.SetToken(*i)
	}
	return pc
}

// SetStorage sets the "storage" field.
func (pc *ProductCreate) SetStorage(i int64) *ProductCreate {
	pc.mutation.SetStorage(i)
	return pc
}

// SetNillableStorage sets the "storage" field if the given value is not nil.
func (pc *ProductCreate) SetNillableStorage(i *int64) *ProductCreate {
	if i != nil {
		pc.SetStorage(*i)
	}
	return pc
}

// Mutation returns the ProductMutation object of the builder.
func (pc *ProductCreate) Mutation() *ProductMutation {
	return pc.mutation
}

// Save creates the Product in the database.
func (pc *ProductCreate) Save(ctx context.Context) (*Product, error) {
	pc.defaults()
	return withHooks[*Product, ProductMutation](ctx, pc.sqlSave, pc.mutation, pc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (pc *ProductCreate) SaveX(ctx context.Context) *Product {
	v, err := pc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (pc *ProductCreate) Exec(ctx context.Context) error {
	_, err := pc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pc *ProductCreate) ExecX(ctx context.Context) {
	if err := pc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (pc *ProductCreate) defaults() {
	if _, ok := pc.mutation.CreatedAt(); !ok {
		v := product.DefaultCreatedAt()
		pc.mutation.SetCreatedAt(v)
	}
	if _, ok := pc.mutation.UpdatedAt(); !ok {
		v := product.DefaultUpdatedAt()
		pc.mutation.SetUpdatedAt(v)
	}
	if _, ok := pc.mutation.Deleted(); !ok {
		v := product.DefaultDeleted
		pc.mutation.SetDeleted(v)
	}
	if _, ok := pc.mutation.Name(); !ok {
		v := product.DefaultName
		pc.mutation.SetName(v)
	}
	if _, ok := pc.mutation.Token(); !ok {
		v := product.DefaultToken
		pc.mutation.SetToken(v)
	}
	if _, ok := pc.mutation.Storage(); !ok {
		v := product.DefaultStorage
		pc.mutation.SetStorage(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (pc *ProductCreate) check() error {
	if _, ok := pc.mutation.Deleted(); !ok {
		return &ValidationError{Name: "deleted", err: errors.New(`ent: missing required field "Product.deleted"`)}
	}
	if _, ok := pc.mutation.UUID(); !ok {
		return &ValidationError{Name: "uuid", err: errors.New(`ent: missing required field "Product.uuid"`)}
	}
	if _, ok := pc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Product.name"`)}
	}
	if _, ok := pc.mutation.Description(); !ok {
		return &ValidationError{Name: "description", err: errors.New(`ent: missing required field "Product.description"`)}
	}
	if _, ok := pc.mutation.Price(); !ok {
		return &ValidationError{Name: "price", err: errors.New(`ent: missing required field "Product.price"`)}
	}
	if _, ok := pc.mutation.Token(); !ok {
		return &ValidationError{Name: "token", err: errors.New(`ent: missing required field "Product.token"`)}
	}
	if _, ok := pc.mutation.Storage(); !ok {
		return &ValidationError{Name: "storage", err: errors.New(`ent: missing required field "Product.storage"`)}
	}
	return nil
}

func (pc *ProductCreate) sqlSave(ctx context.Context) (*Product, error) {
	if err := pc.check(); err != nil {
		return nil, err
	}
	_node, _spec := pc.createSpec()
	if err := sqlgraph.CreateNode(ctx, pc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	pc.mutation.id = &_node.ID
	pc.mutation.done = true
	return _node, nil
}

func (pc *ProductCreate) createSpec() (*Product, *sqlgraph.CreateSpec) {
	var (
		_node = &Product{config: pc.config}
		_spec = sqlgraph.NewCreateSpec(product.Table, sqlgraph.NewFieldSpec(product.FieldID, field.TypeInt))
	)
	_spec.OnConflict = pc.conflict
	if value, ok := pc.mutation.CreatedAt(); ok {
		_spec.SetField(product.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := pc.mutation.UpdatedAt(); ok {
		_spec.SetField(product.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := pc.mutation.Deleted(); ok {
		_spec.SetField(product.FieldDeleted, field.TypeBool, value)
		_node.Deleted = value
	}
	if value, ok := pc.mutation.UUID(); ok {
		_spec.SetField(product.FieldUUID, field.TypeUUID, value)
		_node.UUID = value
	}
	if value, ok := pc.mutation.Name(); ok {
		_spec.SetField(product.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := pc.mutation.Description(); ok {
		_spec.SetField(product.FieldDescription, field.TypeString, value)
		_node.Description = value
	}
	if value, ok := pc.mutation.Price(); ok {
		_spec.SetField(product.FieldPrice, field.TypeFloat64, value)
		_node.Price = value
	}
	if value, ok := pc.mutation.Token(); ok {
		_spec.SetField(product.FieldToken, field.TypeInt64, value)
		_node.Token = value
	}
	if value, ok := pc.mutation.Storage(); ok {
		_spec.SetField(product.FieldStorage, field.TypeInt64, value)
		_node.Storage = value
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Product.Create().
//		SetCreatedAt(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ProductUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (pc *ProductCreate) OnConflict(opts ...sql.ConflictOption) *ProductUpsertOne {
	pc.conflict = opts
	return &ProductUpsertOne{
		create: pc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Product.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (pc *ProductCreate) OnConflictColumns(columns ...string) *ProductUpsertOne {
	pc.conflict = append(pc.conflict, sql.ConflictColumns(columns...))
	return &ProductUpsertOne{
		create: pc,
	}
}

type (
	// ProductUpsertOne is the builder for "upsert"-ing
	//  one Product node.
	ProductUpsertOne struct {
		create *ProductCreate
	}

	// ProductUpsert is the "OnConflict" setter.
	ProductUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdatedAt sets the "updated_at" field.
func (u *ProductUpsert) SetUpdatedAt(v time.Time) *ProductUpsert {
	u.Set(product.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *ProductUpsert) UpdateUpdatedAt() *ProductUpsert {
	u.SetExcluded(product.FieldUpdatedAt)
	return u
}

// SetDeleted sets the "deleted" field.
func (u *ProductUpsert) SetDeleted(v bool) *ProductUpsert {
	u.Set(product.FieldDeleted, v)
	return u
}

// UpdateDeleted sets the "deleted" field to the value that was provided on create.
func (u *ProductUpsert) UpdateDeleted() *ProductUpsert {
	u.SetExcluded(product.FieldDeleted)
	return u
}

// SetName sets the "name" field.
func (u *ProductUpsert) SetName(v string) *ProductUpsert {
	u.Set(product.FieldName, v)
	return u
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *ProductUpsert) UpdateName() *ProductUpsert {
	u.SetExcluded(product.FieldName)
	return u
}

// SetDescription sets the "description" field.
func (u *ProductUpsert) SetDescription(v string) *ProductUpsert {
	u.Set(product.FieldDescription, v)
	return u
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *ProductUpsert) UpdateDescription() *ProductUpsert {
	u.SetExcluded(product.FieldDescription)
	return u
}

// SetPrice sets the "price" field.
func (u *ProductUpsert) SetPrice(v float64) *ProductUpsert {
	u.Set(product.FieldPrice, v)
	return u
}

// UpdatePrice sets the "price" field to the value that was provided on create.
func (u *ProductUpsert) UpdatePrice() *ProductUpsert {
	u.SetExcluded(product.FieldPrice)
	return u
}

// AddPrice adds v to the "price" field.
func (u *ProductUpsert) AddPrice(v float64) *ProductUpsert {
	u.Add(product.FieldPrice, v)
	return u
}

// SetToken sets the "token" field.
func (u *ProductUpsert) SetToken(v int64) *ProductUpsert {
	u.Set(product.FieldToken, v)
	return u
}

// UpdateToken sets the "token" field to the value that was provided on create.
func (u *ProductUpsert) UpdateToken() *ProductUpsert {
	u.SetExcluded(product.FieldToken)
	return u
}

// AddToken adds v to the "token" field.
func (u *ProductUpsert) AddToken(v int64) *ProductUpsert {
	u.Add(product.FieldToken, v)
	return u
}

// SetStorage sets the "storage" field.
func (u *ProductUpsert) SetStorage(v int64) *ProductUpsert {
	u.Set(product.FieldStorage, v)
	return u
}

// UpdateStorage sets the "storage" field to the value that was provided on create.
func (u *ProductUpsert) UpdateStorage() *ProductUpsert {
	u.SetExcluded(product.FieldStorage)
	return u
}

// AddStorage adds v to the "storage" field.
func (u *ProductUpsert) AddStorage(v int64) *ProductUpsert {
	u.Add(product.FieldStorage, v)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.Product.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *ProductUpsertOne) UpdateNewValues() *ProductUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(product.FieldCreatedAt)
		}
		if _, exists := u.create.mutation.UUID(); exists {
			s.SetIgnore(product.FieldUUID)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Product.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *ProductUpsertOne) Ignore() *ProductUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ProductUpsertOne) DoNothing() *ProductUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ProductCreate.OnConflict
// documentation for more info.
func (u *ProductUpsertOne) Update(set func(*ProductUpsert)) *ProductUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ProductUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *ProductUpsertOne) SetUpdatedAt(v time.Time) *ProductUpsertOne {
	return u.Update(func(s *ProductUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *ProductUpsertOne) UpdateUpdatedAt() *ProductUpsertOne {
	return u.Update(func(s *ProductUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeleted sets the "deleted" field.
func (u *ProductUpsertOne) SetDeleted(v bool) *ProductUpsertOne {
	return u.Update(func(s *ProductUpsert) {
		s.SetDeleted(v)
	})
}

// UpdateDeleted sets the "deleted" field to the value that was provided on create.
func (u *ProductUpsertOne) UpdateDeleted() *ProductUpsertOne {
	return u.Update(func(s *ProductUpsert) {
		s.UpdateDeleted()
	})
}

// SetName sets the "name" field.
func (u *ProductUpsertOne) SetName(v string) *ProductUpsertOne {
	return u.Update(func(s *ProductUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *ProductUpsertOne) UpdateName() *ProductUpsertOne {
	return u.Update(func(s *ProductUpsert) {
		s.UpdateName()
	})
}

// SetDescription sets the "description" field.
func (u *ProductUpsertOne) SetDescription(v string) *ProductUpsertOne {
	return u.Update(func(s *ProductUpsert) {
		s.SetDescription(v)
	})
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *ProductUpsertOne) UpdateDescription() *ProductUpsertOne {
	return u.Update(func(s *ProductUpsert) {
		s.UpdateDescription()
	})
}

// SetPrice sets the "price" field.
func (u *ProductUpsertOne) SetPrice(v float64) *ProductUpsertOne {
	return u.Update(func(s *ProductUpsert) {
		s.SetPrice(v)
	})
}

// AddPrice adds v to the "price" field.
func (u *ProductUpsertOne) AddPrice(v float64) *ProductUpsertOne {
	return u.Update(func(s *ProductUpsert) {
		s.AddPrice(v)
	})
}

// UpdatePrice sets the "price" field to the value that was provided on create.
func (u *ProductUpsertOne) UpdatePrice() *ProductUpsertOne {
	return u.Update(func(s *ProductUpsert) {
		s.UpdatePrice()
	})
}

// SetToken sets the "token" field.
func (u *ProductUpsertOne) SetToken(v int64) *ProductUpsertOne {
	return u.Update(func(s *ProductUpsert) {
		s.SetToken(v)
	})
}

// AddToken adds v to the "token" field.
func (u *ProductUpsertOne) AddToken(v int64) *ProductUpsertOne {
	return u.Update(func(s *ProductUpsert) {
		s.AddToken(v)
	})
}

// UpdateToken sets the "token" field to the value that was provided on create.
func (u *ProductUpsertOne) UpdateToken() *ProductUpsertOne {
	return u.Update(func(s *ProductUpsert) {
		s.UpdateToken()
	})
}

// SetStorage sets the "storage" field.
func (u *ProductUpsertOne) SetStorage(v int64) *ProductUpsertOne {
	return u.Update(func(s *ProductUpsert) {
		s.SetStorage(v)
	})
}

// AddStorage adds v to the "storage" field.
func (u *ProductUpsertOne) AddStorage(v int64) *ProductUpsertOne {
	return u.Update(func(s *ProductUpsert) {
		s.AddStorage(v)
	})
}

// UpdateStorage sets the "storage" field to the value that was provided on create.
func (u *ProductUpsertOne) UpdateStorage() *ProductUpsertOne {
	return u.Update(func(s *ProductUpsert) {
		s.UpdateStorage()
	})
}

// Exec executes the query.
func (u *ProductUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for ProductCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ProductUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *ProductUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *ProductUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// ProductCreateBulk is the builder for creating many Product entities in bulk.
type ProductCreateBulk struct {
	config
	builders []*ProductCreate
	conflict []sql.ConflictOption
}

// Save creates the Product entities in the database.
func (pcb *ProductCreateBulk) Save(ctx context.Context) ([]*Product, error) {
	specs := make([]*sqlgraph.CreateSpec, len(pcb.builders))
	nodes := make([]*Product, len(pcb.builders))
	mutators := make([]Mutator, len(pcb.builders))
	for i := range pcb.builders {
		func(i int, root context.Context) {
			builder := pcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ProductMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, pcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = pcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, pcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, pcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (pcb *ProductCreateBulk) SaveX(ctx context.Context) []*Product {
	v, err := pcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (pcb *ProductCreateBulk) Exec(ctx context.Context) error {
	_, err := pcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pcb *ProductCreateBulk) ExecX(ctx context.Context) {
	if err := pcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Product.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.ProductUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (pcb *ProductCreateBulk) OnConflict(opts ...sql.ConflictOption) *ProductUpsertBulk {
	pcb.conflict = opts
	return &ProductUpsertBulk{
		create: pcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Product.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (pcb *ProductCreateBulk) OnConflictColumns(columns ...string) *ProductUpsertBulk {
	pcb.conflict = append(pcb.conflict, sql.ConflictColumns(columns...))
	return &ProductUpsertBulk{
		create: pcb,
	}
}

// ProductUpsertBulk is the builder for "upsert"-ing
// a bulk of Product nodes.
type ProductUpsertBulk struct {
	create *ProductCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Product.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *ProductUpsertBulk) UpdateNewValues() *ProductUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(product.FieldCreatedAt)
			}
			if _, exists := b.mutation.UUID(); exists {
				s.SetIgnore(product.FieldUUID)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Product.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *ProductUpsertBulk) Ignore() *ProductUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *ProductUpsertBulk) DoNothing() *ProductUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the ProductCreateBulk.OnConflict
// documentation for more info.
func (u *ProductUpsertBulk) Update(set func(*ProductUpsert)) *ProductUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&ProductUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *ProductUpsertBulk) SetUpdatedAt(v time.Time) *ProductUpsertBulk {
	return u.Update(func(s *ProductUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *ProductUpsertBulk) UpdateUpdatedAt() *ProductUpsertBulk {
	return u.Update(func(s *ProductUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetDeleted sets the "deleted" field.
func (u *ProductUpsertBulk) SetDeleted(v bool) *ProductUpsertBulk {
	return u.Update(func(s *ProductUpsert) {
		s.SetDeleted(v)
	})
}

// UpdateDeleted sets the "deleted" field to the value that was provided on create.
func (u *ProductUpsertBulk) UpdateDeleted() *ProductUpsertBulk {
	return u.Update(func(s *ProductUpsert) {
		s.UpdateDeleted()
	})
}

// SetName sets the "name" field.
func (u *ProductUpsertBulk) SetName(v string) *ProductUpsertBulk {
	return u.Update(func(s *ProductUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *ProductUpsertBulk) UpdateName() *ProductUpsertBulk {
	return u.Update(func(s *ProductUpsert) {
		s.UpdateName()
	})
}

// SetDescription sets the "description" field.
func (u *ProductUpsertBulk) SetDescription(v string) *ProductUpsertBulk {
	return u.Update(func(s *ProductUpsert) {
		s.SetDescription(v)
	})
}

// UpdateDescription sets the "description" field to the value that was provided on create.
func (u *ProductUpsertBulk) UpdateDescription() *ProductUpsertBulk {
	return u.Update(func(s *ProductUpsert) {
		s.UpdateDescription()
	})
}

// SetPrice sets the "price" field.
func (u *ProductUpsertBulk) SetPrice(v float64) *ProductUpsertBulk {
	return u.Update(func(s *ProductUpsert) {
		s.SetPrice(v)
	})
}

// AddPrice adds v to the "price" field.
func (u *ProductUpsertBulk) AddPrice(v float64) *ProductUpsertBulk {
	return u.Update(func(s *ProductUpsert) {
		s.AddPrice(v)
	})
}

// UpdatePrice sets the "price" field to the value that was provided on create.
func (u *ProductUpsertBulk) UpdatePrice() *ProductUpsertBulk {
	return u.Update(func(s *ProductUpsert) {
		s.UpdatePrice()
	})
}

// SetToken sets the "token" field.
func (u *ProductUpsertBulk) SetToken(v int64) *ProductUpsertBulk {
	return u.Update(func(s *ProductUpsert) {
		s.SetToken(v)
	})
}

// AddToken adds v to the "token" field.
func (u *ProductUpsertBulk) AddToken(v int64) *ProductUpsertBulk {
	return u.Update(func(s *ProductUpsert) {
		s.AddToken(v)
	})
}

// UpdateToken sets the "token" field to the value that was provided on create.
func (u *ProductUpsertBulk) UpdateToken() *ProductUpsertBulk {
	return u.Update(func(s *ProductUpsert) {
		s.UpdateToken()
	})
}

// SetStorage sets the "storage" field.
func (u *ProductUpsertBulk) SetStorage(v int64) *ProductUpsertBulk {
	return u.Update(func(s *ProductUpsert) {
		s.SetStorage(v)
	})
}

// AddStorage adds v to the "storage" field.
func (u *ProductUpsertBulk) AddStorage(v int64) *ProductUpsertBulk {
	return u.Update(func(s *ProductUpsert) {
		s.AddStorage(v)
	})
}

// UpdateStorage sets the "storage" field to the value that was provided on create.
func (u *ProductUpsertBulk) UpdateStorage() *ProductUpsertBulk {
	return u.Update(func(s *ProductUpsert) {
		s.UpdateStorage()
	})
}

// Exec executes the query.
func (u *ProductUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the ProductCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for ProductCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *ProductUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
