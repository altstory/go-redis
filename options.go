package redis

import (
	"time"

	"github.com/go-redis/redis"
)

// SetOption 代表设置一个 key 时候用到的各种选项。
type SetOption struct {
	t      setOptionType
	expire time.Duration
}

// Expire 返回一个用于设置 key 超时的选项。
// 详见 https://redis.io/commands/set。
func Expire(timeout time.Duration) SetOption {
	return SetOption{
		t:      setOptionExpire,
		expire: timeout,
	}
}

// XX 返回一个用于设置仅当 key 存在才改变其值的选项。
// 详见 https://redis.io/commands/set。
func XX() SetOption {
	return SetOption{
		t: setOptionXX,
	}
}

// NX 返回一个用于设置仅当 key 不存在才改变其值的选项。
// 详见 https://redis.io/commands/set。
func NX() SetOption {
	return SetOption{
		t: setOptionNX,
	}
}

// Args 返回用于拼接 Redis 命令的参数。
func (so *SetOption) Args() []interface{} {
	switch so.t {
	case setOptionExpire:
		if so.expire%time.Second == 0 {
			return []interface{}{"EX", int64(so.expire / time.Second)}
		}

		return []interface{}{"PX", int64(so.expire.Round(time.Millisecond) / time.Millisecond)}

	case setOptionNX:
		return []interface{}{"NX"}
	case setOptionXX:
		return []interface{}{"XX"}
	}

	return nil
}

type setOptionType int

const (
	setOptionInvalid setOptionType = iota
	setOptionExpire
	setOptionNX
	setOptionXX
)

// StoreOption 代表 sorted set 对计算结果进行存储时的选项。
type StoreOption struct {
	t   storeOptionType
	opt interface{}
}

// Weights 返回一个用于 sorted set 计算结果存储的选项 WEIGHTS。
// 详见 https://redis.io/commands/zunionstore。
func Weights(weights ...float64) StoreOption {
	return StoreOption{
		t:   storeOptionWeights,
		opt: weights,
	}
}

// AggregateSum 返回一个用于 sorted set 计算结果存储的选项 AGGREGATE SUM。
// 详见 https://redis.io/commands/zunionstore。
func AggregateSum() StoreOption {
	return StoreOption{
		t:   storeOptionAggregate,
		opt: "SUM",
	}
}

// AggregateMin 返回一个用于 sorted set 计算结果存储的选项 AGGREGATE MIN。
// 详见 https://redis.io/commands/zunionstore。
func AggregateMin() StoreOption {
	return StoreOption{
		t:   storeOptionAggregate,
		opt: "MIN",
	}
}

// AggregateMax 返回一个用于 sorted set 计算结果存储的选项 AGGREGATE MAX。
// 详见 https://redis.io/commands/zunionstore。
func AggregateMax() StoreOption {
	return StoreOption{
		t:   storeOptionAggregate,
		opt: "MAX",
	}
}

// Args 返回用于拼接 Redis 命令的参数。
func (so *StoreOption) Args() []interface{} {
	switch so.t {
	case storeOptionWeights:
		weights := so.opt.([]float64)

		if len(weights) == 0 {
			return nil
		}

		args := make([]interface{}, 0, len(weights)+1)
		args = append(args, "WEIGHTS")

		for _, w := range weights {
			args = append(args, w)
		}

		return args

	case storeOptionAggregate:
		return []interface{}{"AGGREGATE", so.opt}
	}

	return nil
}

func (so *StoreOption) fillZStore(zstore *redis.ZStore) {
	switch so.t {
	case storeOptionWeights:
		weights := so.opt.([]float64)
		zstore.Weights = weights

	case storeOptionAggregate:
		aggregate := so.opt.(string)
		zstore.Aggregate = aggregate
	}
}

type storeOptionType int

const (
	storeOptionInvalid storeOptionType = iota
	storeOptionWeights
	storeOptionAggregate
)

// RangeOption 代表 sorted set 对计算结果进行范围查询时的选项。
type RangeOption struct {
	offset int
	count  int
}

// Limit 返回一个用于 sorted set 范围查询的 limit。
func Limit(offset int, count int) RangeOption {
	return RangeOption{
		offset: offset,
		count:  count,
	}
}

// Args 返回用于拼接 Redis 命令的参数。
func (ro *RangeOption) Args() []interface{} {
	return []interface{}{"LIMIT", ro.offset, ro.count}
}

func (ro *RangeOption) fillZRangeBy(zrangeby *redis.ZRangeBy) {
	zrangeby.Offset = int64(ro.offset)
	zrangeby.Count = int64(ro.count)
}

// ScanOption 代表一个扫描选项。
type ScanOption struct {
	t   scanOptionType
	opt interface{}
}

// Match 返回一个扫描选项，用于在各种 Scan 中实现 MATCH pattern。
func Match(pattern string) ScanOption {
	return ScanOption{
		t:   scanOptionMatch,
		opt: pattern,
	}
}

// Count 返回一个扫描选项，用于在各种 Scan 中实现 COUNT count。
func Count(count int) ScanOption {
	return ScanOption{
		t:   scanOptionCount,
		opt: count,
	}
}

// Args 返回用于拼接 Redis 命令的参数。
func (so *ScanOption) Args() []interface{} {
	switch so.t {
	case scanOptionMatch:
		return []interface{}{"MATCH", so.opt}
	case scanOptionCount:
		return []interface{}{"COUNT", so.opt}
	}

	return nil
}

type scanOptionType int

const (
	scanOptionTypeInvalid scanOptionType = iota
	scanOptionMatch
	scanOptionCount
)

// FlushOption 代表 FLUSHDB 和 FLUSHALL 的选项。
type FlushOption int

const (
	flushOptionInvalid FlushOption = iota
	flushOptionAsync
)

// Async 代表用异步模式进行 flush。
func Async() FlushOption {
	return flushOptionAsync
}
