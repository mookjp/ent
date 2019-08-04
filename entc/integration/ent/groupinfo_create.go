// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"fbc/ent/entc/integration/ent/group"
	"fbc/ent/entc/integration/ent/groupinfo"

	"fbc/ent/dialect"
	"fbc/ent/dialect/gremlin"
	"fbc/ent/dialect/gremlin/graph/dsl"
	"fbc/ent/dialect/gremlin/graph/dsl/__"
	"fbc/ent/dialect/gremlin/graph/dsl/g"
	"fbc/ent/dialect/gremlin/graph/dsl/p"
	"fbc/ent/dialect/sql"
)

// GroupInfoCreate is the builder for creating a GroupInfo entity.
type GroupInfoCreate struct {
	config
	desc      *string
	max_users *int
	groups    map[string]struct{}
}

// SetDesc sets the desc field.
func (gic *GroupInfoCreate) SetDesc(s string) *GroupInfoCreate {
	gic.desc = &s
	return gic
}

// SetMaxUsers sets the max_users field.
func (gic *GroupInfoCreate) SetMaxUsers(i int) *GroupInfoCreate {
	gic.max_users = &i
	return gic
}

// SetNillableMaxUsers sets the max_users field if the given value is not nil.
func (gic *GroupInfoCreate) SetNillableMaxUsers(i *int) *GroupInfoCreate {
	if i != nil {
		gic.SetMaxUsers(*i)
	}
	return gic
}

// AddGroupIDs adds the groups edge to Group by ids.
func (gic *GroupInfoCreate) AddGroupIDs(ids ...string) *GroupInfoCreate {
	if gic.groups == nil {
		gic.groups = make(map[string]struct{})
	}
	for i := range ids {
		gic.groups[ids[i]] = struct{}{}
	}
	return gic
}

// AddGroups adds the groups edges to Group.
func (gic *GroupInfoCreate) AddGroups(g ...*Group) *GroupInfoCreate {
	ids := make([]string, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return gic.AddGroupIDs(ids...)
}

// Save creates the GroupInfo in the database.
func (gic *GroupInfoCreate) Save(ctx context.Context) (*GroupInfo, error) {
	if gic.desc == nil {
		return nil, errors.New("ent: missing required field \"desc\"")
	}
	if gic.max_users == nil {
		v := groupinfo.DefaultMaxUsers
		gic.max_users = &v
	}
	switch gic.driver.Dialect() {
	case dialect.MySQL, dialect.SQLite:
		return gic.sqlSave(ctx)
	case dialect.Neptune:
		return gic.gremlinSave(ctx)
	default:
		return nil, errors.New("ent: unsupported dialect")
	}
}

// SaveX calls Save and panics if Save returns an error.
func (gic *GroupInfoCreate) SaveX(ctx context.Context) *GroupInfo {
	v, err := gic.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (gic *GroupInfoCreate) sqlSave(ctx context.Context) (*GroupInfo, error) {
	var (
		res sql.Result
		gi  = &GroupInfo{config: gic.config}
	)
	tx, err := gic.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	builder := sql.Insert(groupinfo.Table).Default(gic.driver.Dialect())
	if gic.desc != nil {
		builder.Set(groupinfo.FieldDesc, *gic.desc)
		gi.Desc = *gic.desc
	}
	if gic.max_users != nil {
		builder.Set(groupinfo.FieldMaxUsers, *gic.max_users)
		gi.MaxUsers = *gic.max_users
	}
	query, args := builder.Query()
	if err := tx.Exec(ctx, query, args, &res); err != nil {
		return nil, rollback(tx, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, rollback(tx, err)
	}
	gi.ID = strconv.FormatInt(id, 10)
	if len(gic.groups) > 0 {
		p := sql.P()
		for eid := range gic.groups {
			eid, err := strconv.Atoi(eid)
			if err != nil {
				return nil, rollback(tx, err)
			}
			p.Or().EQ(group.FieldID, eid)
		}
		query, args := sql.Update(groupinfo.GroupsTable).
			Set(groupinfo.GroupsColumn, id).
			Where(sql.And(p, sql.IsNull(groupinfo.GroupsColumn))).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(gic.groups) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"groups\" %v already connected to a different \"GroupInfo\"", keys(gic.groups))})
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return gi, nil
}

func (gic *GroupInfoCreate) gremlinSave(ctx context.Context) (*GroupInfo, error) {
	res := &gremlin.Response{}
	query, bindings := gic.gremlin().Query()
	if err := gic.driver.Exec(ctx, query, bindings, res); err != nil {
		return nil, err
	}
	if err, ok := isConstantError(res); ok {
		return nil, err
	}
	gi := &GroupInfo{config: gic.config}
	if err := gi.FromResponse(res); err != nil {
		return nil, err
	}
	return gi, nil
}

func (gic *GroupInfoCreate) gremlin() *dsl.Traversal {
	type constraint struct {
		pred *dsl.Traversal // constraint predicate.
		test *dsl.Traversal // test matches and its constant.
	}
	constraints := make([]*constraint, 0, 1)
	v := g.AddV(groupinfo.Label)
	if gic.desc != nil {
		v.Property(dsl.Single, groupinfo.FieldDesc, *gic.desc)
	}
	if gic.max_users != nil {
		v.Property(dsl.Single, groupinfo.FieldMaxUsers, *gic.max_users)
	}
	for id := range gic.groups {
		v.AddE(group.InfoLabel).From(g.V(id)).InV()
		constraints = append(constraints, &constraint{
			pred: g.E().HasLabel(group.InfoLabel).OutV().HasID(id).Count(),
			test: __.Is(p.NEQ(0)).Constant(NewErrUniqueEdge(groupinfo.Label, group.InfoLabel, id)),
		})
	}
	if len(constraints) == 0 {
		return v.ValueMap(true)
	}
	tr := constraints[0].pred.Coalesce(constraints[0].test, v.ValueMap(true))
	for _, cr := range constraints[1:] {
		tr = cr.pred.Coalesce(cr.test, tr)
	}
	return tr
}