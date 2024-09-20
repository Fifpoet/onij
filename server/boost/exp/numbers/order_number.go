package numbers

import (
	"onij/boost/collection/collext"
	"onij/boost/conv"
	"onij/boost/exp"
	"strconv"
	"strings"
)

type OrderNumber struct {
	v    string
	nums []uint64
}

func NewOrderNumber(number string) *OrderNumber {
	n := &OrderNumber{v: number}
	parts := strings.Split(number, dot)
	nums := make([]uint64, 0, len(parts))
	for _, v := range parts {
		if num, err := strconv.ParseUint(v, 10, 64); err == nil {
			nums = append(nums, num)
			continue
		}
		nums = nil
		break
	}
	n.nums = nums
	return n
}

func (r *OrderNumber) Sub(order uint) *OrderNumber {
	n := new(OrderNumber)
	if len(r.nums) != 0 {
		n.nums = append(n.nums, r.nums...)
	}
	n.nums = append(n.nums, uint64(order))
	n.v = strings.Join(collext.Pick(n.nums, conv.UInt64ToString), dot)
	return n
}

func (r *OrderNumber) String() string {
	if r == nil {
		return exp.Zero[string]()
	}
	return r.v
}

func (r *OrderNumber) Compare(t *OrderNumber) int {
	if r == nil && t == nil {
		return 0
	}
	if r == nil && t != nil {
		return -1
	}
	if r != nil && t == nil {
		return 1
	}
	for i := 0; i < max(len(r.nums), len(t.nums)); i++ {
		rNum, _ := collext.Index(r.nums, i)
		tNum, _ := collext.Index(t.nums, i)
		if rNum < tNum {
			return -1
		}
		if rNum > tNum {
			return 1
		}
	}
	return 0
}

func (r *OrderNumber) Equals(t *OrderNumber) bool { return r.Compare(t) == 0 }

func (r *OrderNumber) Copy() *OrderNumber {
	if r == nil {
		return nil
	}
	return &OrderNumber{
		v:    r.v,
		nums: append(r.nums[:0:0], r.nums...),
	}
}
