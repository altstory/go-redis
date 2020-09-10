package redis

import (
	"github.com/altstory/go-redis/internal/driver"
	"github.com/go-redis/redis"
)

// SortedSets 代表 Redis 跟 zset 相关的接口，详见 https://redis.io/commands#sorted_set。
type SortedSets interface {
	// TODO: BZPOPMAX
	// TODO: BZPOPMIN

	ZAdd(key string, mss ...MemberAndScore) (added int, err error)
	ZCard(key string) (count int, err error)
	ZCount(key string, min ScoreRange, max ScoreRange) (count int, err error)
	ZIncrBy(key string, incr float64, member string) (score float64, err error)
	ZInterStore(dst string, keys []string, options ...StoreOption) (count int, err error)
	ZLexCount(key string, min MemberRange, max MemberRange) (count int, err error)
	ZPopMax(key string) (ms MemberAndScore, err error)
	ZPopMaxN(key string, count int) (mss MemberAndScores, err error)
	ZPopMin(key string) (ms MemberAndScore, err error)
	ZPopMinN(key string, count int) (mss MemberAndScores, err error)
	ZRange(key string, start float64, stop float64) (members []BulkString, err error)
	ZRangeWithScores(key string, start float64, stop float64) (mss MemberAndScores, err error)
	ZRangeByLex(key string, min MemberRange, max MemberRange, options ...RangeOption) (members []BulkString, err error)
	ZRangeByScore(key string, min ScoreRange, max ScoreRange, options ...RangeOption) (members []BulkString, err error)
	ZRangeByScoreWithScores(key string, min ScoreRange, max ScoreRange, options ...RangeOption) (mss MemberAndScores, err error)
	ZRank(key string, member string) (rank int, exists bool, err error)
	ZRem(key string, members ...string) (removed int, err error)
	ZRemRangeByLex(key string, min MemberRange, max MemberRange) (removed int, err error)
	ZRemRangeByRank(key string, start int, stop int) (removed int, err error)
	ZRemRangeByScore(key string, min ScoreRange, max ScoreRange) (removed int, err error)
	ZRevRange(key string, start float64, stop float64) (members []BulkString, err error)
	ZRevRangeWithScores(key string, start float64, stop float64) (mss MemberAndScores, err error)
	ZRevRangeByLex(key string, min MemberRange, max MemberRange, options ...RangeOption) (members []BulkString, err error)
	ZRevRangeByScore(key string, min ScoreRange, max ScoreRange, options ...RangeOption) (members []BulkString, err error)
	ZRevRangeByScoreWithScores(key string, min ScoreRange, max ScoreRange, options ...RangeOption) (mss MemberAndScores, err error)
	ZRevRank(key string, member string) (rank int, exists bool, err error)
	ZScore(key string, member string) (score float64, exists bool, err error)
	ZUnionStore(dst string, keys []string, options ...StoreOption) (count int, err error)
}

func (r *redisImpl) ZAdd(key string, mss ...MemberAndScore) (added int, err error) {
	if len(mss) == 0 {
		return
	}

	zs := make([]redis.Z, 0, len(mss))

	for _, ms := range mss {
		zs = append(zs, redis.Z{
			Member: ms.Member,
			Score:  ms.Score,
		})
	}

	err = r.do("ZADD", func(client driver.Client) error {
		added, err = mustBeInt(client, client.ZAdd(key, zs...))
		return err
	})
	return
}

func (r *redisImpl) ZCard(key string) (count int, err error) {
	err = r.do("ZCARD", func(client driver.Client) error {
		count, err = mustBeInt(client, client.ZCard(key))
		return err
	})
	return
}

func (r *redisImpl) ZCount(key string, min ScoreRange, max ScoreRange) (count int, err error) {
	err = r.do("ZCOUNT", func(client driver.Client) error {
		count, err = mustBeInt(client, client.ZCount(key, min.String(), max.String()))
		return err
	})
	return
}

func (r *redisImpl) ZIncrBy(key string, incr float64, member string) (score float64, err error) {
	err = r.do("ZINCRBY", func(client driver.Client) error {
		score, err = mustBeFloat64(client, client.ZIncrBy(key, incr, member))
		return err
	})
	return
}

func (r *redisImpl) ZInterStore(dst string, keys []string, options ...StoreOption) (count int, err error) {
	if len(keys) == 0 {
		return
	}

	var zstore redis.ZStore

	for _, opt := range options {
		opt.fillZStore(&zstore)
	}

	err = r.do("ZINTERSTORE", func(client driver.Client) error {
		count, err = mustBeInt(client, client.ZInterStore(dst, zstore, keys...))
		return err
	})
	return
}

func (r *redisImpl) ZLexCount(key string, min MemberRange, max MemberRange) (count int, err error) {
	err = r.do("ZLEXCOUNT", func(client driver.Client) error {
		count, err = mustBeInt(client, client.ZLexCount(key, min.String(), max.String()))
		return err
	})
	return
}

func (r *redisImpl) ZPopMax(key string) (ms MemberAndScore, err error) {
	err = r.do("ZPOPMAX", func(client driver.Client) error {
		ms, err = mustBeMemberAndScore(client, client.ZPopMax(key))
		return err
	})
	return
}

func (r *redisImpl) ZPopMaxN(key string, count int) (mss MemberAndScores, err error) {
	err = r.do("ZPOPMAX-N", func(client driver.Client) error {
		mss, err = mustBeMemberAndScores(client, client.ZPopMax(key, int64(count)))
		return err
	})
	return
}

func (r *redisImpl) ZPopMin(key string) (ms MemberAndScore, err error) {
	err = r.do("ZPOPMIN", func(client driver.Client) error {
		ms, err = mustBeMemberAndScore(client, client.ZPopMin(key))
		return err
	})
	return
}

func (r *redisImpl) ZPopMinN(key string, count int) (mss MemberAndScores, err error) {
	err = r.do("ZPOPMIN-N", func(client driver.Client) error {
		mss, err = mustBeMemberAndScores(client, client.ZPopMin(key, int64(count)))
		return err
	})
	return
}

func (r *redisImpl) ZRange(key string, start float64, stop float64) (members []BulkString, err error) {
	err = r.do("ZRANGE", func(client driver.Client) error {
		cmd := redis.NewStringSliceCmd("ZRANGE", key, start, stop)

		if err = client.Process(cmd); err != nil {
			return err
		}

		members, err = mustBeBulkStrings(client, cmd)
		return err
	})
	return
}

func (r *redisImpl) ZRangeWithScores(key string, start float64, stop float64) (mss MemberAndScores, err error) {
	err = r.do("ZRANGE-WITHSCORES", func(client driver.Client) error {
		cmd := redis.NewSliceCmd("ZRANGE", key, start, stop, "WITHSCORES")

		if err = client.Process(cmd); err != nil {
			return err
		}

		mss, err = mustBeMemberAndScores(client, cmd)
		return err
	})
	return
}

func (r *redisImpl) ZRangeByLex(key string, min MemberRange, max MemberRange, options ...RangeOption) (members []BulkString, err error) {
	zrangeby := redis.ZRangeBy{
		Min: min.String(),
		Max: max.String(),
	}

	for _, opt := range options {
		opt.fillZRangeBy(&zrangeby)
	}

	err = r.do("ZRANGEBYLEX", func(client driver.Client) error {
		members, err = mustBeBulkStrings(client, client.ZRangeByLex(key, zrangeby))
		return err
	})
	return
}

func (r *redisImpl) ZRangeByScore(key string, min ScoreRange, max ScoreRange, options ...RangeOption) (members []BulkString, err error) {
	zrangeby := redis.ZRangeBy{
		Min: min.String(),
		Max: max.String(),
	}

	for _, opt := range options {
		opt.fillZRangeBy(&zrangeby)
	}

	err = r.do("ZRANGEBYSCORE", func(client driver.Client) error {
		members, err = mustBeBulkStrings(client, client.ZRangeByScore(key, zrangeby))
		return err
	})
	return
}

func (r *redisImpl) ZRangeByScoreWithScores(key string, min ScoreRange, max ScoreRange, options ...RangeOption) (mss MemberAndScores, err error) {
	zrangeby := redis.ZRangeBy{
		Min: min.String(),
		Max: max.String(),
	}

	for _, opt := range options {
		opt.fillZRangeBy(&zrangeby)
	}

	err = r.do("ZRANGEBYSCORE-WITHSCORES", func(client driver.Client) error {
		mss, err = mustBeMemberAndScores(client, client.ZRangeByScoreWithScores(key, zrangeby))
		return err
	})
	return
}

func (r *redisImpl) ZRank(key string, member string) (rank int, exists bool, err error) {
	err = r.do("ZRANK", func(client driver.Client) error {
		cmd := redis.NewIntCmd("ZRANK", key, member)

		if err = client.Process(cmd); err != nil {
			if err != redis.Nil {
				return err
			}

			return nil
		}

		exists = true
		rank, err = mustBeInt(client, cmd)
		return err
	})
	return
}

func (r *redisImpl) ZRem(key string, members ...string) (removed int, err error) {
	if len(members) == 0 {
		return
	}

	data := make([]interface{}, 0, len(members))

	for _, m := range members {
		data = append(data, m)
	}

	err = r.do("ZREM", func(client driver.Client) error {
		removed, err = mustBeInt(client, client.ZRem(key, data...))
		return err
	})
	return
}

func (r *redisImpl) ZRemRangeByLex(key string, min MemberRange, max MemberRange) (removed int, err error) {
	err = r.do("ZREMRANGEBYLEX", func(client driver.Client) error {
		removed, err = mustBeInt(client, client.ZRemRangeByLex(key, min.String(), max.String()))
		return err
	})
	return
}

func (r *redisImpl) ZRemRangeByRank(key string, start int, stop int) (removed int, err error) {
	err = r.do("ZREMRANGEBYRANK", func(client driver.Client) error {
		removed, err = mustBeInt(client, client.ZRemRangeByRank(key, int64(start), int64(stop)))
		return err
	})
	return
}

func (r *redisImpl) ZRemRangeByScore(key string, min ScoreRange, max ScoreRange) (removed int, err error) {
	err = r.do("ZREMRANGEBYSCORE", func(client driver.Client) error {
		removed, err = mustBeInt(client, client.ZRemRangeByScore(key, min.String(), max.String()))
		return err
	})
	return
}

func (r *redisImpl) ZRevRange(key string, start float64, stop float64) (members []BulkString, err error) {
	err = r.do("ZREVRANGE", func(client driver.Client) error {
		cmd := redis.NewStringSliceCmd("ZREVRANGE", key, start, stop)

		if err = client.Process(cmd); err != nil {
			return err
		}

		members, err = mustBeBulkStrings(client, cmd)
		return err
	})
	return
}

func (r *redisImpl) ZRevRangeWithScores(key string, start float64, stop float64) (mss MemberAndScores, err error) {
	err = r.do("ZREVRANGE-WITHSCORES", func(client driver.Client) error {
		cmd := redis.NewSliceCmd("ZREVRANGE", key, start, stop, "WITHSCORES")

		if err = client.Process(cmd); err != nil {
			return err
		}

		mss, err = mustBeMemberAndScores(client, cmd)
		return err
	})
	return
}

func (r *redisImpl) ZRevRangeByLex(key string, min MemberRange, max MemberRange, options ...RangeOption) (members []BulkString, err error) {
	zrangeby := redis.ZRangeBy{
		Min: min.String(),
		Max: max.String(),
	}

	for _, opt := range options {
		opt.fillZRangeBy(&zrangeby)
	}

	err = r.do("ZREVRANGEBYLEX", func(client driver.Client) error {
		members, err = mustBeBulkStrings(client, client.ZRevRangeByLex(key, zrangeby))
		return err
	})
	return
}

func (r *redisImpl) ZRevRangeByScore(key string, min ScoreRange, max ScoreRange, options ...RangeOption) (members []BulkString, err error) {
	zrangeby := redis.ZRangeBy{
		Min: min.String(),
		Max: max.String(),
	}

	for _, opt := range options {
		opt.fillZRangeBy(&zrangeby)
	}

	err = r.do("ZREVRANGEBYSCORE", func(client driver.Client) error {
		members, err = mustBeBulkStrings(client, client.ZRevRangeByScore(key, zrangeby))
		return err
	})
	return
}

func (r *redisImpl) ZRevRangeByScoreWithScores(key string, min ScoreRange, max ScoreRange, options ...RangeOption) (mss MemberAndScores, err error) {
	zrangeby := redis.ZRangeBy{
		Min: min.String(),
		Max: max.String(),
	}

	for _, opt := range options {
		opt.fillZRangeBy(&zrangeby)
	}

	err = r.do("ZREVRANGEBYSCORE-WITHSCORES", func(client driver.Client) error {
		mss, err = mustBeMemberAndScores(client, client.ZRevRangeByScoreWithScores(key, zrangeby))
		return err
	})
	return
}

func (r *redisImpl) ZRevRank(key string, member string) (rank int, exists bool, err error) {
	err = r.do("ZREVRANK", func(client driver.Client) error {
		cmd := redis.NewIntCmd("ZREVRANK", key, member)

		if err = client.Process(cmd); err != nil {
			if err != redis.Nil {
				return err
			}

			return nil
		}

		exists = true
		rank, err = mustBeInt(client, cmd)
		return err
	})
	return
}

func (r *redisImpl) ZScore(key string, member string) (score float64, exists bool, err error) {
	err = r.do("ZSCORE", func(client driver.Client) error {
		cmd := redis.NewFloatCmd("ZSCORE", key, member)

		if err = client.Process(cmd); err != nil {
			if err != redis.Nil {
				return err
			}

			return nil
		}

		exists = true
		score, err = mustBeFloat64(client, cmd)
		return err
	})
	return
}

func (r *redisImpl) ZUnionStore(dst string, keys []string, options ...StoreOption) (count int, err error) {
	if len(keys) == 0 {
		return
	}

	var zstore redis.ZStore

	for _, opt := range options {
		opt.fillZStore(&zstore)
	}

	err = r.do("ZUNIONSTORE", func(client driver.Client) error {
		count, err = mustBeInt(client, client.ZUnionStore(dst, zstore, keys...))
		return err
	})
	return
}
