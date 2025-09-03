package cdb

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type Proxy interface {
	DB() *gorm.DB
	With(tx *gorm.DB) Proxy
	R(ctx context.Context, opts ...Option) *gorm.DB
	W(ctx context.Context, opts ...Option) *gorm.DB
}

type DefaultProxy struct {
	db *gorm.DB
}

type ShardConfig struct {
	Table    string `json:"table,omitempty"`
	ShardNum int32  `json:"shard_num,omitempty"`
	ShardKey string `json:"shard_key,omitempty"`
}

// NewDefaultProxy 创建默认的数据库代理。
func NewDefaultProxy(db *gorm.DB) *DefaultProxy {
	return &DefaultProxy{db: db}
}

// DB 返回 gorm.DB 对象。
func (p *DefaultProxy) DB() *gorm.DB {
	return p.db
}

// With 返回一个新的 Proxy，使用指定的事务连接。
func (p *DefaultProxy) With(tx *gorm.DB) Proxy {
	return &DefaultProxy{
		db: tx,
	}
}

func (p *DefaultProxy) W(ctx context.Context, opts ...Option) *gorm.DB {
	conn := p.db.WithContext(ctx).Clauses(dbresolver.Write)
	for _, opt := range opts {
		conn = opt(conn)
	}
	return conn
}

func (p *DefaultProxy) R(ctx context.Context, opts ...Option) *gorm.DB {
	conn := p.db.WithContext(ctx).Clauses(dbresolver.Read)
	for _, opt := range opts {
		conn = opt(conn)
	}
	return conn
}
