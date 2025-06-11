package util

import "context"

func IfElse[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

type Page struct {
	Page  int32
	Limit int32
}

func (p *Page) OffsetNum() int {
	return (int(p.Page) - 1) * int(p.Limit)
}
func (p *Page) LimitNum() int {
	return int(p.Limit)
}
func (p *Page) PageNum() int {
	return int(p.Page)
}

func PageGetAll[R any](ctx context.Context, limit int32, maxCircle int32, f func(ctx context.Context, offset, limit int32) (page []R, stop bool, err error)) (res []R, err error) {
	offset := int32(0)
	for range maxCircle {
		pageRes, stop, err := f(ctx, offset, limit)
		if err != nil {
			return nil, err
		}
		res = append(res, pageRes...)
		if len(pageRes) < int(limit) {
			return res, nil
		}
		if stop {
			return res, nil
		}
		offset += limit
	}
	return
}
