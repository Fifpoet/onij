package cdb

import (
	"reflect"
	"strings"

	"gorm.io/gorm"
)

const (
	OpType_Eq    = "="
	OpType_Neq   = "!="
	OpType_Gt    = ">"
	OpType_Gte   = ">="
	OpType_Lt    = "<"
	OpType_Lte   = "<="
	OpType_Btw   = "BETWEEN"
	OpType_In    = "IN"
	OpType_NotIn = "NOT IN"
	OpType_Like  = "LIKE"

	InnerQuerierKey_Asc     = "__[inner_asc]__"
	InnerQuerierKey_Desc    = "__[inner_desc]__"
	InnerQuerierKey_Limit   = "__[inner_limit]__"
	InnerQuerierKey_Unscope = "__[inner_unscope]__"
)

type Option func(*gorm.DB) *gorm.DB

func NewOption(col string, val any) Option {
	if op, ok := val.(*Op); ok {
		return op.ToOption(col)
	}
	valType := reflect.TypeOf(val)
	if valType.Kind() == reflect.Slice {
		return IN(val).ToOption(col)
	}
	return EQ(val).ToOption(col)
}

// WithAsc 字段增序
// Deprecated: 计划废弃，请使用相应的 Querier 中的 WithAsc 替代
func WithAsc(cols ...string) Option {
	return func(db *gorm.DB) *gorm.DB {
		if len(cols) == 0 {
			return db
		}
		return db.Order(strings.Join(cols, ","))
	}
}

// WithDesc 字段倒序
// Deprecated: 计划废弃，请使用相应的 Querier 中的 WithDesc 替代
func WithDesc(cols ...string) Option {
	return func(db *gorm.DB) *gorm.DB {
		if len(cols) == 0 {
			return db
		}
		return db.Order(strings.Join(cols, " DESC,") + " DESC")
	}
}

func EQ(val any) *Op {
	return newOp(OpType_Eq, val)
}

func NEQ(val any) *Op {
	return newOp(OpType_Neq, val)
}

func GT(val any) *Op {
	return newOp(OpType_Gt, val)
}

func GTE(val any) *Op {
	return newOp(OpType_Gte, val)
}

func LT(val any) *Op {
	return newOp(OpType_Lt, val)
}

func LTE(val any) *Op {
	return newOp(OpType_Lte, val)
}

func BTW(val1, val2 any) *Op {
	return newOp(OpType_Btw, val1, val2)
}

func IN(val any) *Op {
	return newOp(OpType_In, val)
}

func NOTIN(val any) *Op {
	return newOp(OpType_NotIn, val)
}

func LIKE(val any) *Op {
	return newOp(OpType_Like, val)
}

type Op struct {
	Values []any
	Type   string
}

func newOp(t string, values ...any) *Op {
	return &Op{Values: values, Type: t}
}

func (op *Op) ToOption(col string) Option {
	switch col {
	case InnerQuerierKey_Asc:
		return func(db *gorm.DB) *gorm.DB {
			for _, orderCol := range op.Values {
				for _, col := range orderCol.([]string) {
					db = db.Order(col)
				}
			}
			return db
		}
	case InnerQuerierKey_Desc:
		return func(db *gorm.DB) *gorm.DB {
			for _, orderCol := range op.Values {
				for _, col := range orderCol.([]string) {
					db = db.Order(col + " DESC")
				}
			}
			return db
		}
	case InnerQuerierKey_Limit:
		return func(db *gorm.DB) *gorm.DB {
			return db.Limit(op.Values[0].(int))
		}
	case InnerQuerierKey_Unscope:
		return func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}
	default:
		switch op.Type {
		case OpType_Eq:
			return func(db *gorm.DB) *gorm.DB {
				return db.Where(col+" = ?", op.Values[0])
			}
		case OpType_Neq:
			return func(db *gorm.DB) *gorm.DB {
				return db.Where(col+" != ?", op.Values[0])
			}
		case OpType_Gt:
			return func(db *gorm.DB) *gorm.DB {
				return db.Where(col+" > ?", op.Values[0])
			}
		case OpType_Gte:
			return func(db *gorm.DB) *gorm.DB {
				return db.Where(col+" >= ?", op.Values[0])
			}
		case OpType_Lt:
			return func(db *gorm.DB) *gorm.DB {
				return db.Where(col+" < ?", op.Values[0])
			}
		case OpType_Lte:
			return func(db *gorm.DB) *gorm.DB {
				return db.Where(col+" <= ?", op.Values[0])
			}
		case OpType_Btw:
			return func(db *gorm.DB) *gorm.DB {
				return db.Where(col+" BETWEEN ? AND ?", op.Values[0], op.Values[1])
			}
		case OpType_In:
			return func(db *gorm.DB) *gorm.DB {
				return db.Where(col+" IN (?)", op.Values...)
			}
		case OpType_NotIn:
			return func(db *gorm.DB) *gorm.DB {
				return db.Where(col+" NOT IN (?)", op.Values...)
			}
		case OpType_Like:
			return func(db *gorm.DB) *gorm.DB {
				return db.Where(col+" LIKE ?", op.Values[0])
			}
		default:
			return nil
		}
	}
}

type BasicQuerier map[string]any

func (q BasicQuerier) Add(col string, val any) {
	q[col] = val
}

func (q BasicQuerier) WithAsc(cols ...string) {
	q[InnerQuerierKey_Asc] = cols
}

func (q BasicQuerier) WithDesc(cols ...string) {
	q[InnerQuerierKey_Desc] = cols
}

func (q BasicQuerier) WithLimit(limit int) {
	q[InnerQuerierKey_Limit] = limit
}

func (q BasicQuerier) WithUnscoped() {
	q[InnerQuerierKey_Unscope] = true
}

func (q BasicQuerier) ToOptions() []Option {
	opts := make([]Option, 0, len(q))
	for col, val := range q {
		opts = append(opts, NewOption(col, val))
	}
	clear(q)
	return opts
}

type BasicUpdater map[string]any

func (u BasicUpdater) Add(col string, val any) {
	u[col] = val
}

func (u BasicUpdater) ToMap() map[string]any {
	if v, ok := u["deleted_at"]; ok {
		if reflect.ValueOf(v).IsZero() {
			u["deleted_at"] = nil
		}
	}
	return u
}

func (u BasicUpdater) IsEmpty() bool {
	return len(u) == 0
}
