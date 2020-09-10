package redis

import (
	"time"

	"github.com/altstory/go-redis/internal/driver"
	"github.com/go-redis/redis"
)

// Strings 代表 Redis 跟字符串相关的接口，详见 https://redis.io/commands#string。
type Strings interface {
	Append(key string, value string) (l int, err error)

	// TODO: BitCount(key string) (count int, err error)
	// TODO: BitCountWithRange(key string, start int, end int) (count int, err error)
	// TODO: BitField
	// TODO: BitOp
	// TODO: BitPos

	Decr(key string) (value int64, err error)
	DecrBy(key string, decr int64) (value int64, err error)
	Get(key string) (value BulkString, err error)

	// TODO: GetBit

	GetRange(key string, start int, end int) (value BulkString, err error)
	GetSet(key string, value string) (old BulkString, err error)
	Incr(key string) (value int64, err error)
	IncrBy(key string, incr int64) (value int64, err error)
	IncrByFloat(key string, incr float64) (value float64, err error)
	MGet(keys ...string) (values []BulkString, err error)
	MSet(kvs ...KeyAndValue) (err error)
	MSetNX(kvs ...KeyAndValue) (isSet bool, err error)
	Set(key string, value string, options ...SetOption) (isSet bool, err error)

	// TODO: SetBit

	SetEx(key string, timeout time.Duration, value string) (err error)
	SetNX(key string, value string) (isSet bool, err error)
	SetRange(key string, offset int, value string) (modified int, err error)
	StrLen(key string) (l int, err error)
}

func (r *redisImpl) Append(key string, value string) (l int, err error) {
	err = r.do("APPEND", func(client driver.Client) error {
		l, err = mustBeInt(client, client.Append(key, value))
		return err
	})
	return
}

func (r *redisImpl) Decr(key string) (value int64, err error) {
	err = r.do("DECR", func(client driver.Client) error {
		value, err = mustBeInt64(client, client.Decr(key))
		return err
	})
	return
}

func (r *redisImpl) DecrBy(key string, decr int64) (value int64, err error) {
	err = r.do("DECRBY", func(client driver.Client) error {
		value, err = mustBeInt64(client, client.DecrBy(key, decr))
		return err
	})
	return
}

func (r *redisImpl) Get(key string) (value BulkString, err error) {
	err = r.do("GET", func(client driver.Client) error {
		value, err = mustBeBulkString(client, client.Get(key))
		return err
	})
	return
}

func (r *redisImpl) GetRange(key string, start int, end int) (value BulkString, err error) {
	err = r.do("GETRANGE", func(client driver.Client) error {
		value, err = mustBeBulkString(client, client.GetRange(key, int64(start), int64(end)))
		return err
	})
	return
}

func (r *redisImpl) GetSet(key string, value string) (old BulkString, err error) {
	err = r.do("GETSET", func(client driver.Client) error {
		old, err = mustBeBulkString(client, client.GetSet(key, value))
		return err
	})
	return
}

func (r *redisImpl) Incr(key string) (value int64, err error) {
	err = r.do("INCR", func(client driver.Client) error {
		value, err = mustBeInt64(client, client.Incr(key))
		return err
	})
	return
}

func (r *redisImpl) IncrBy(key string, incr int64) (value int64, err error) {
	err = r.do("INCRBY", func(client driver.Client) error {
		value, err = mustBeInt64(client, client.IncrBy(key, incr))
		return err
	})
	return
}

func (r *redisImpl) IncrByFloat(key string, incr float64) (value float64, err error) {
	err = r.do("INCRBYFLOAT", func(client driver.Client) error {
		value, err = mustBeFloat64(client, client.IncrByFloat(key, incr))
		return err
	})
	return
}

func (r *redisImpl) MGet(keys ...string) (values []BulkString, err error) {
	if len(keys) == 0 {
		return
	}

	err = r.do("MGET", func(client driver.Client) error {
		values, err = mustBeBulkStrings(client, client.MGet(keys...))
		return err
	})
	return
}

func (r *redisImpl) MSet(kvs ...KeyAndValue) (err error) {
	if len(kvs) == 0 {
		return
	}

	values := make([]interface{}, 0, 2*len(kvs))

	for _, kv := range kvs {
		values = append(values, kv.Key, kv.Value)
	}

	err = r.do("MSET", func(client driver.Client) error {
		_, err = mustBeStatus(client, client.MSet(values...))
		return err
	})
	return
}

func (r *redisImpl) MSetNX(kvs ...KeyAndValue) (isSet bool, err error) {
	if len(kvs) == 0 {
		return
	}

	values := make([]interface{}, 0, 2*len(kvs))

	for _, kv := range kvs {
		values = append(values, kv.Key, kv.Value)
	}

	err = r.do("MSETNX", func(client driver.Client) error {
		isSet, err = mustBeBool(client, client.MSetNX(values...))
		return err
	})
	return
}

func (r *redisImpl) Set(key string, value string, options ...SetOption) (isSet bool, err error) {
	err = r.do("SET", func(client driver.Client) error {
		args := make([]interface{}, 0, 8) // SET 最多有这么多参数。
		args = append(args, "SET", key, value)

		for _, opt := range options {
			args = append(args, opt.Args()...)
		}

		cmd := redis.NewStatusCmd(args...)

		if err = client.Process(cmd); err != nil {
			if err != redis.Nil {
				return err
			}

			return nil
		}

		isSet = true
		_, err = mustBeStatus(client, cmd)
		return err
	})
	return
}

func (r *redisImpl) SetEx(key string, timeout time.Duration, value string) (err error) {
	err = r.do("SETEX", func(client driver.Client) error {
		_, err = mustBeStatus(client, client.Set(key, value, timeout))
		return err
	})
	return
}

func (r *redisImpl) SetNX(key string, value string) (isSet bool, err error) {
	err = r.do("SETNX", func(client driver.Client) error {
		isSet, err = mustBeBool(client, client.SetNX(key, value, 0))
		return err
	})
	return
}

func (r *redisImpl) SetRange(key string, offset int, value string) (modified int, err error) {
	err = r.do("SETRANGE", func(client driver.Client) error {
		modified, err = mustBeInt(client, client.SetRange(key, int64(offset), value))
		return err
	})
	return
}

func (r *redisImpl) StrLen(key string) (l int, err error) {
	err = r.do("STRLEN", func(client driver.Client) error {
		l, err = mustBeInt(client, client.StrLen(key))
		return err
	})
	return
}
