package uid

import (
	"math"
	"math/rand"
	"sync/atomic"
	"time"
)

const (
	bitsTotal     = 32                       // 总位数
	bitsHight     = 4                        // 高位部分位数
	bitsTimestamp = 8                        // 时间戳部分位数
	bitsIncr      = 4                        // 递增部分位数
	bitsLow       = bitsTimestamp + bitsIncr // 低位部分位数
)

// Generator 生成 uid(unique identifiers)
type Generator interface {
	Generate() int64
}

type Options struct {
	Since time.Time // 生成的起始时间
	Group uint16    // id 所属分组
}

type Option func(*Options)

func WithSince(since time.Time) Option {
	return func(o *Options) {
		o.Since = since
	}
}

func WithGroup(group uint16) Option {
	return func(o *Options) {
		o.Group = group
	}
}

func WithDefault() Option {
	return func(o *Options) {
		o.Since = time.Now()
		o.Group = uint16(rand.Intn(math.MaxUint16))
	}
}

func NewGenerator(opts ...Option) Generator {
	g := &generator{o: new(Options)}
	for _, opt := range append([]Option{WithDefault()}, opts...) {
		opt(g.o)
	}
	g.init()
	return g
}

// generator 根据 Options 配置生成 uid
//
// High(16 bits): Group
// Low(48 bits): Timestamp(40 bits) + Incr(8 bits)
type generator struct {
	o    *Options // 配置
	high uint64   // uid 高位
	low  uint64   // uid 低位
}

// Generate 生成一个 uid
func (r *generator) Generate() int64 { return r.generate() }

// init 根据配置初始化生成器
func (r *generator) init() {
	r.high = uint64(r.o.Group) << bitsLow
	r.low = r.intercept(uint64(r.o.Since.UnixMilli()), bitsTimestamp) << bitsIncr
}

// generate 高低位计算生成 uid
func (r *generator) generate() int64 {
	return int64(r.high | r.intercept(atomic.AddUint64(&r.low, 1), bitsLow))
}

// intercept 获取将 s 从低位起截取 n 位后的结果
func (r *generator) intercept(s uint64, n uint) uint64 {
	return s & (math.MaxUint64 >> (bitsTotal - n))
}
