package exp

const (
	DefaultPagingLimit    = 10 // 默认分页大小
	DefaultMaxPagingLimit = 50 // 最大分页大小
)

type Paging interface {
	Offset() int
	Limit() int
}

type PagingOption func(*paging)

func WithPagingOptionLimit(limit int) PagingOption {
	return func(o *paging) { o.limit = limit }
}

func WithPagingOptionMaxLimit(max int) PagingOption {
	return func(o *paging) { o.maxLimit = max }
}

// DefaultPagingLimit = 10, DefaultMaxPagingLimit = 50
func NewPaging(page, size int, opts ...PagingOption) Paging {
	paging := &paging{page: page, size: size}
	for _, opt := range opts {
		opt(paging)
	}
	paging.normalize()
	return paging
}

type paging struct {
	page int
	size int

	limit    int
	maxLimit int
}

func (r *paging) Offset() int {
	n := r.page
	if n <= 0 {
		n = 1
	}
	return (n - 1) * r.Limit()
}

func (r *paging) Limit() int {
	if r.size <= 0 {
		return r.limit
	}
	if r.size > r.maxLimit {
		return r.maxLimit
	}
	return r.size
}

func (r *paging) normalize() {
	if r.limit <= 0 {
		r.limit = DefaultPagingLimit
	}
	if r.maxLimit <= 0 {
		r.maxLimit = DefaultMaxPagingLimit
	}
}
