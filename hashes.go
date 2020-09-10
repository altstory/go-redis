package redis

import (
	"github.com/altstory/go-redis/internal/driver"
)

// Hashes 代表 Redis 跟 hash 相关的接口，详见 https://redis.io/commands#hash。
type Hashes interface {
	HDel(key string, fields ...string) (deleted int, err error)
	HExists(key string, field string) (exists bool, err error)
	HGet(key string, field string) (value BulkString, err error)
	HGetAll(key string) (fieldAndValues KeyAndValues, err error)
	HIncrBy(key string, field string, incr int64) (value int64, err error)
	HIncrByFloat(key string, field string, incr float64) (value float64, err error)
	HKeys(key string) (keys []BulkString, err error)
	HLen(key string) (l int, err error)
	HMGet(key string, fields ...string) (values []BulkString, err error)
	HMSet(key string, fieldAndValues ...KeyAndValue) (err error)
	HSet(key string, field string, value string) (isNew bool, err error)
	HSetNX(key string, field string, value string) (isNew bool, err error)

	// HStrLen(key string, field string) (l int, err error) // TODO: driver 没提供，先跳过

	HVals(key string) (values []BulkString, err error)
}

func (r *redisImpl) HDel(key string, fields ...string) (deleted int, err error) {
	if len(fields) == 0 {
		return
	}

	err = r.do("HDEL", func(client driver.Client) error {
		deleted, err = mustBeInt(client, client.HDel(key, fields...))
		return err
	})
	return
}

func (r *redisImpl) HExists(key string, field string) (exists bool, err error) {
	err = r.do("HEXISTS", func(client driver.Client) error {
		exists, err = mustBeBool(client, client.HExists(key, field))
		return err
	})
	return
}

func (r *redisImpl) HGet(key string, field string) (value BulkString, err error) {
	err = r.do("HGET", func(client driver.Client) error {
		value, err = mustBeBulkString(client, client.HGet(key, field))
		return err
	})
	return
}

func (r *redisImpl) HGetAll(key string) (fieldAndValues KeyAndValues, err error) {
	err = r.do("HGETALL", func(client driver.Client) error {
		fieldAndValues, err = mustBeKeyAndValues(client, client.HGetAll(key))
		return err
	})
	return
}

func (r *redisImpl) HIncrBy(key string, field string, incr int64) (value int64, err error) {
	err = r.do("HINCRBY", func(client driver.Client) error {
		value, err = mustBeInt64(client, client.HIncrBy(key, field, incr))
		return err
	})
	return
}

func (r *redisImpl) HIncrByFloat(key string, field string, incr float64) (value float64, err error) {
	err = r.do("HINCRBYFLOAT", func(client driver.Client) error {
		value, err = mustBeFloat64(client, client.HIncrByFloat(key, field, incr))
		return err
	})
	return
}

func (r *redisImpl) HKeys(key string) (keys []BulkString, err error) {
	err = r.do("HKEYS", func(client driver.Client) error {
		keys, err = mustBeBulkStrings(client, client.HKeys(key))
		return err
	})
	return
}

func (r *redisImpl) HLen(key string) (l int, err error) {
	err = r.do("HLEN", func(client driver.Client) error {
		l, err = mustBeInt(client, client.HLen(key))
		return err
	})
	return
}

func (r *redisImpl) HMGet(key string, fields ...string) (values []BulkString, err error) {
	if len(fields) == 0 {
		return
	}

	err = r.do("HMGET", func(client driver.Client) error {
		values, err = mustBeBulkStrings(client, client.HMGet(key, fields...))
		return err
	})
	return
}

func (r *redisImpl) HMSet(key string, fieldAndValues ...KeyAndValue) (err error) {
	if len(fieldAndValues) == 0 {
		return
	}

	fields := make(map[string]interface{}, len(fieldAndValues))

	for _, fv := range fieldAndValues {
		fields[fv.Key] = fv.Value
	}

	err = r.do("HMSET", func(client driver.Client) error {
		_, err = mustBeStatus(client, client.HMSet(key, fields))
		return err
	})
	return
}

func (r *redisImpl) HSet(key string, field string, value string) (isNew bool, err error) {
	err = r.do("HSET", func(client driver.Client) error {
		isNew, err = mustBeBool(client, client.HSet(key, field, value))
		return err
	})
	return
}

func (r *redisImpl) HSetNX(key string, field string, value string) (isNew bool, err error) {
	err = r.do("HSETNX", func(client driver.Client) error {
		isNew, err = mustBeBool(client, client.HSetNX(key, field, value))
		return err
	})
	return
}

// func (r *redisImpl) HStrLen(key string, field string) (l int, err error) {
// 	err = r.do("HSTRLEN", func(client driver.Client) error {
// 		l, err = mustBeInt(client, client.HStrLen(key))
// 		return err
// 	})
// 	return
// }

func (r *redisImpl) HVals(key string) (values []BulkString, err error) {
	err = r.do("HVALS", func(client driver.Client) error {
		values, err = mustBeBulkStrings(client, client.HVals(key))
		return err
	})
	return
}
