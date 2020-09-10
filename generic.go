package redis

import (
	"time"

	"github.com/altstory/go-redis/internal/driver"
)

// Generic 代表 Redis 各种经典 K/V 接口，详见 https://redis.io/commands#generic。
//
// 注意，WAIT 不是用户使用的命令，这里不支持。
type Generic interface {
	Del(keys ...string) (deleted int, err error)
	Dump(key string) (value BulkString, err error)
	Exists(keys ...string) (existing int, err error)
	Expire(key string, timeout time.Duration) (isSet bool, err error)
	ExpireAt(key string, t time.Time) (isSet bool, err error)
	Keys(pattern string) (keys []BulkString, err error)

	// TODO: Migrate
	// TODO: Move
	// TODO: Object

	Persist(key string) (persisted bool, err error)
	RandomKey() (key BulkString, err error)
	Rename(old, new string) (err error)
	RenameNX(old, new string) (renamed bool, err error)

	// TODO: Restore
	// TODO: Store

	Touch(keys ...string) (touched int, err error)
	TTL(key string) (ttl time.Duration, err error)
	Type(key string) (keyType KeyType, err error)
	Unlink(keys ...string) (unlinked int, err error)
}

// KeyType 代表 Redis key 所对应的类型。
type KeyType string

// 所有 Redis key 类型。
const (
	TypeString KeyType = "string"
	TypeList   KeyType = "list"
	TypeSet    KeyType = "set"
	TypeZSet   KeyType = "zset"
	TypeHash   KeyType = "hash"
	TypeStream KeyType = "stream"
)

func (r *redisImpl) Del(keys ...string) (deleted int, err error) {
	if len(keys) == 0 {
		return
	}

	err = r.do("DEL", func(client driver.Client) error {
		deleted, err = mustBeInt(client, client.Del(keys...))
		return err
	})
	return
}

func (r *redisImpl) Dump(key string) (value BulkString, err error) {
	err = r.do("DUMP", func(client driver.Client) error {
		value, err = mustBeBulkString(client, client.Dump(key))
		return err
	})
	return
}

func (r *redisImpl) Exists(keys ...string) (existing int, err error) {
	if len(keys) == 0 {
		return
	}

	err = r.do("EXISTS", func(client driver.Client) error {
		existing, err = mustBeInt(client, client.Exists(keys...))
		return err
	})
	return
}

func (r *redisImpl) Expire(key string, timeout time.Duration) (isSet bool, err error) {
	err = r.do("EXPIRE", func(client driver.Client) error {
		isSet, err = mustBeBool(client, client.Expire(key, timeout))
		return err
	})
	return
}

func (r *redisImpl) ExpireAt(key string, t time.Time) (isSet bool, err error) {
	err = r.do("EXPIREAT", func(client driver.Client) error {
		isSet, err = mustBeBool(client, client.ExpireAt(key, t))
		return err
	})
	return
}

func (r *redisImpl) Keys(pattern string) (keys []BulkString, err error) {
	err = r.do("KEYS", func(client driver.Client) error {
		keys, err = mustBeBulkStrings(client, client.Keys(pattern))
		return err
	})
	return
}

func (r *redisImpl) Persist(key string) (persisted bool, err error) {
	err = r.do("PERSIST", func(client driver.Client) error {
		persisted, err = mustBeBool(client, client.Persist(key))
		return err
	})
	return
}

func (r *redisImpl) RandomKey() (key BulkString, err error) {
	err = r.do("RANDOMKEY", func(client driver.Client) error {
		key, err = mustBeBulkString(client, client.RandomKey())
		return err
	})
	return
}

func (r *redisImpl) Rename(old, new string) (err error) {
	err = r.do("RENAME", func(client driver.Client) error {
		_, err = mustBeStatus(client, client.Rename(old, new))
		return err
	})
	return
}

func (r *redisImpl) RenameNX(old, new string) (renamed bool, err error) {
	err = r.do("RENAMENX", func(client driver.Client) error {
		renamed, err = mustBeBool(client, client.RenameNX(old, new))
		return err
	})
	return
}

func (r *redisImpl) Touch(keys ...string) (touched int, err error) {
	if len(keys) == 0 {
		return
	}

	err = r.do("TOUCH", func(client driver.Client) error {
		touched, err = mustBeInt(client, client.Touch(keys...))
		return err
	})
	return
}

func (r *redisImpl) TTL(key string) (ttl time.Duration, err error) {
	err = r.do("TTL", func(client driver.Client) error {
		ttl, err = mustBeDuration(client, client.TTL(key))
		return err
	})
	return
}

func (r *redisImpl) Type(key string) (keyType KeyType, err error) {
	err = r.do("TYPE", func(client driver.Client) error {
		var t string
		t, err = mustBeStatus(client, client.Type(key))
		keyType = KeyType(t)
		return err
	})
	return
}

func (r *redisImpl) Unlink(keys ...string) (unlinked int, err error) {
	if len(keys) == 0 {
		return
	}

	err = r.do("UNLINK", func(client driver.Client) error {
		unlinked, err = mustBeInt(client, client.Unlink(keys...))
		return err
	})
	return
}
