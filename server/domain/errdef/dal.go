package errdef

import "errors"

// entity

var (
	ErrDalEntityIdInvalid              = errors.New("invalid model id")         // 无效实体Id
	ErrDalEntityUpdatableColumnInvalid = errors.New("invalid updatable column") // 无效实体可更新列
)

// update

var (
	ErrDalUpdateColumnsUnmatched = errors.New("update columns unmatched") // 更新列不匹配
)
