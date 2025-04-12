package util

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
