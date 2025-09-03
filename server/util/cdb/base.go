package cdb

import (
	"context"
	"database/sql"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const maxLimit = 500

var (
	Err_NoLimit         = fmt.Errorf("no limit")
	Err_EmptyConditions = fmt.Errorf("empty conditions")
	Err_EmptyUpdateInfo = fmt.Errorf("empty update info")
	Err_TooManyRows     = fmt.Errorf("too many rows")
)

type Interface[Dal any] interface {
	With(tx *gorm.DB) Dal
	Tx(ctx context.Context, fn func(*gorm.DB) error, opts ...*sql.TxOptions) error
	R(ctx context.Context, opts ...Option) *gorm.DB
	W(ctx context.Context, opts ...Option) *gorm.DB
}

type TableModel interface {
	TableName() string
	UpdatableColumns() []string
}

type Querier interface {
	~map[string]any
	ToOptions() []Option
}

type Updater interface {
	~map[string]any
	ToMap() map[string]any
	IsEmpty() bool
}

type Dal[M TableModel, Q Querier, U Updater] struct {
	p    Proxy
	with bool
}

func NewDal[M TableModel, Q Querier, U Updater](proxy Proxy) *Dal[M, Q, U] {
	return &Dal[M, Q, U]{p: proxy}
}

func (d *Dal[M, Q, U]) R(ctx context.Context, opts ...Option) *gorm.DB {
	var conn *gorm.DB
	if d.with {
		conn = d.p.DB()
	} else {
		conn = d.p.R(ctx, opts...)
	}
	return conn
}

func (d *Dal[M, Q, U]) W(ctx context.Context, opts ...Option) *gorm.DB {
	var conn *gorm.DB
	if d.with {
		conn = d.p.DB()
	} else {
		conn = d.p.W(ctx)
	}
	for _, opt := range opts {
		conn = opt(conn)
	}
	return conn
}

func (d *Dal[M, Q, U]) With(tx *gorm.DB) *Dal[M, Q, U] {
	return &Dal[M, Q, U]{p: d.p.With(tx), with: true}
}

func (d *Dal[M, Q, U]) Tx(ctx context.Context, fn func(*gorm.DB) error, opts ...*sql.TxOptions) error {
	return d.W(ctx).Transaction(fn, opts...)
}

func (d *Dal[M, Q, U]) Q() Q {
	return make(Q)
}

func (d *Dal[M, Q, U]) U() U {
	return make(U)
}

func (d *Dal[M, Q, U]) QueryFirst(ctx context.Context, opts ...Option) (*M, error) {
	if len(opts) == 0 {
		return nil, Err_EmptyConditions
	}
	var res *M
	err := d.R(ctx, opts...).First(&res).Error
	return res, err
}

func (d *Dal[M, Q, U]) QuerySlice(ctx context.Context, offset, limit int, opts ...Option) ([]*M, error) {
	if len(opts) == 0 {
		return nil, Err_EmptyConditions
	}
	if limit > maxLimit {
		return nil, Err_TooManyRows
	}
	if limit <= 0 {
		return nil, Err_NoLimit
	}

	var res []*M
	err := d.R(ctx, opts...).Offset(offset).Limit(limit).Find(&res).Error
	return res, err
}

func (d *Dal[M, Q, U]) QueryAll(ctx context.Context, opts ...Option) ([]*M, error) {
	// TODO: 翻页
	if len(opts) == 0 {
		return nil, Err_EmptyConditions
	}

	var res []*M
	err := d.R(ctx, opts...).Find(&res).Error
	return res, err
}

func (d *Dal[M, Q, U]) Count(ctx context.Context, opts ...Option) (int64, error) {
	if len(opts) == 0 {
		return 0, Err_EmptyConditions
	}

	var count int64
	err := d.R(ctx, opts...).Model(new(M)).Count(&count).Error
	return count, err
}

func (d *Dal[M, Q, U]) Create(ctx context.Context, mds ...*M) (int64, error) {
	if len(mds) == 0 {
		return 0, nil
	}

	conn := d.W(ctx).CreateInBatches(mds, maxLimit)
	return conn.RowsAffected, conn.Error
}

func (d *Dal[M, Q, U]) Update(ctx context.Context, u U, conditions ...Option) (int64, error) {
	if len(conditions) == 0 {
		return 0, Err_EmptyConditions
	}
	if u.IsEmpty() {
		return 0, Err_EmptyUpdateInfo
	}

	conn := d.W(ctx, conditions...).Model(new(M)).Updates(u.ToMap())
	return conn.RowsAffected, conn.Error
}

func (d *Dal[M, Q, U]) Upsert(ctx context.Context, data ...*M) (int64, error) {
	if len(data) == 0 {
		return 0, nil
	}

	var m M
	//goland:noinspection GoDfaNilDereference
	cols := m.UpdatableColumns()
	conflictClause := clause.OnConflict{UpdateAll: true}
	if len(cols) != 0 {
		conflictClause = clause.OnConflict{DoUpdates: clause.AssignmentColumns(cols)}
	}

	conn := d.W(ctx).Clauses(conflictClause).CreateInBatches(data, maxLimit)
	return conn.RowsAffected, conn.Error
}

func (d *Dal[M, Q, U]) Delete(ctx context.Context, opts ...Option) (int64, error) {
	if len(opts) == 0 {
		return 0, Err_EmptyConditions
	}

	conn := d.W(ctx, opts...).Delete(new(M))
	return conn.RowsAffected, conn.Error
}
