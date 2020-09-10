package redis

import (
	"github.com/altstory/go-redis/internal/driver"
)

// Lists 代表 Redis 跟 list 相关的接口，详见 https://redis.io/commands#list。
type Lists interface {
	// TODO: BLPop
	// TODO: BRPop
	// TODO: BRPopLPush

	LIndex(key string, index int) (value BulkString, err error)
	LInsertBefore(key string, pivot string, value string) (l int, err error)
	LInsertAfter(key string, pivot string, value string) (l int, err error)
	LLen(key string) (l int, err error)
	LPop(key string) (value BulkString, err error)
	LPush(key string, values ...string) (l int, err error)
	LPushX(key string, value string) (l int, err error)
	LRange(key string, start int, stop int) (values []BulkString, err error)
	LRem(key string, count int, value string) (removed int, err error)
	LSet(key string, index int, value string) (err error)
	LTrim(key string, start int, stop int) (err error)
	RPop(key string) (value BulkString, err error)
	RPopLPush(src string, dst string) (value BulkString, err error)
	RPush(key string, values ...string) (l int, err error)
	RPushX(key string, value string) (l int, err error)
}

func (r *redisImpl) LIndex(key string, index int) (value BulkString, err error) {
	err = r.do("LINDEX", func(client driver.Client) error {
		value, err = mustBeBulkString(client, client.LIndex(key, int64(index)))
		return err
	})
	return
}

func (r *redisImpl) LInsertBefore(key string, pivot string, value string) (l int, err error) {
	err = r.do("LINSERTBEFORE", func(client driver.Client) error {
		l, err = mustBeInt(client, client.LInsertBefore(key, pivot, value))
		return err
	})
	return
}

func (r *redisImpl) LInsertAfter(key string, pivot string, value string) (l int, err error) {
	err = r.do("LINSERTAFTER", func(client driver.Client) error {
		l, err = mustBeInt(client, client.LInsertAfter(key, pivot, value))
		return err
	})
	return
}

func (r *redisImpl) LLen(key string) (l int, err error) {
	err = r.do("LLEN", func(client driver.Client) error {
		l, err = mustBeInt(client, client.HLen(key))
		return err
	})
	return
}

func (r *redisImpl) LPop(key string) (value BulkString, err error) {
	err = r.do("LPOP", func(client driver.Client) error {
		value, err = mustBeBulkString(client, client.LPop(key))
		return err
	})
	return
}

func (r *redisImpl) LPush(key string, values ...string) (l int, err error) {
	if len(values) == 0 {
		return
	}

	data := make([]interface{}, 0, len(values))

	for _, v := range values {
		data = append(data, v)
	}

	err = r.do("LPUSH", func(client driver.Client) error {
		l, err = mustBeInt(client, client.LPush(key, data...))
		return err
	})
	return
}

func (r *redisImpl) LPushX(key string, value string) (l int, err error) {
	err = r.do("LPUSH", func(client driver.Client) error {
		l, err = mustBeInt(client, client.LPush(key, value))
		return err
	})
	return
}

func (r *redisImpl) LRange(key string, start int, stop int) (values []BulkString, err error) {
	err = r.do("LRANGE", func(client driver.Client) error {
		values, err = mustBeBulkStrings(client, client.LRange(key, int64(start), int64(stop)))
		return err
	})
	return
}

func (r *redisImpl) LRem(key string, count int, value string) (removed int, err error) {
	err = r.do("LREM", func(client driver.Client) error {
		removed, err = mustBeInt(client, client.LRem(key, int64(count), value))
		return err
	})
	return
}

func (r *redisImpl) LSet(key string, index int, value string) (err error) {
	err = r.do("LSET", func(client driver.Client) error {
		_, err = mustBeStatus(client, client.LSet(key, int64(index), value))
		return err
	})
	return
}

func (r *redisImpl) LTrim(key string, start int, stop int) (err error) {
	err = r.do("LTRIM", func(client driver.Client) error {
		_, err = mustBeStatus(client, client.LTrim(key, int64(start), int64(stop)))
		return err
	})
	return
}

func (r *redisImpl) RPop(key string) (value BulkString, err error) {
	err = r.do("RPOP", func(client driver.Client) error {
		value, err = mustBeBulkString(client, client.RPop(key))
		return err
	})
	return
}

func (r *redisImpl) RPopLPush(src string, dst string) (value BulkString, err error) {
	err = r.do("RPOPLPUSH", func(client driver.Client) error {
		value, err = mustBeBulkString(client, client.RPopLPush(src, dst))
		return err
	})
	return
}

func (r *redisImpl) RPush(key string, values ...string) (l int, err error) {
	if len(values) == 0 {
		return
	}

	data := make([]interface{}, 0, len(values))

	for _, v := range values {
		data = append(data, v)
	}

	err = r.do("RPUSH", func(client driver.Client) error {
		l, err = mustBeInt(client, client.RPush(key, data...))
		return err
	})
	return
}

func (r *redisImpl) RPushX(key string, value string) (l int, err error) {
	err = r.do("RPUSHX", func(client driver.Client) error {
		l, err = mustBeInt(client, client.RPushX(key, value))
		return err
	})
	return
}
