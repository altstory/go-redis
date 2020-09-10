package redis

import (
	"github.com/altstory/go-redis/internal/driver"
)

// Sets 代表 Redis 跟 set 相关的接口，详见 https://redis.io/commands#set。
type Sets interface {
	SAdd(key string, members ...string) (added int, err error)
	SCard(key string) (count int, err error)
	SDiff(keys ...string) (members []BulkString, err error)
	SDiffStore(dst string, keys ...string) (count int, err error)
	SInter(keys ...string) (members []BulkString, err error)
	SInterStore(dst string, keys ...string) (count int, err error)
	SIsMember(key string, member string) (exists bool, err error)
	SMembers(key string) (members []BulkString, err error)
	SMove(src string, dst string, member string) (moved bool, err error)
	SPop(key string) (member BulkString, err error)
	SPopN(key string, count int) (members []BulkString, err error) // SPOP key count
	SRandMember(key string) (member BulkString, err error)
	SRandMemberN(key string, count int) (members []BulkString, err error) // SRANDMEMBER key count
	SRem(key string, members ...string) (removed int, err error)
	SUnion(keys ...string) (members []BulkString, err error)
	SUnionStore(dst string, keys ...string) (count int, err error)
}

func (r *redisImpl) SAdd(key string, members ...string) (added int, err error) {
	if len(members) == 0 {
		return
	}

	data := make([]interface{}, 0, len(members))

	for _, m := range members {
		data = append(data, m)
	}

	err = r.do("SADD", func(client driver.Client) error {
		added, err = mustBeInt(client, client.SAdd(key, data...))
		return err
	})
	return
}

func (r *redisImpl) SCard(key string) (count int, err error) {
	err = r.do("SCARD", func(client driver.Client) error {
		count, err = mustBeInt(client, client.SCard(key))
		return err
	})
	return
}

func (r *redisImpl) SDiff(keys ...string) (members []BulkString, err error) {
	if len(keys) == 0 {
		return
	}

	err = r.do("SDIFF", func(client driver.Client) error {
		members, err = mustBeBulkStrings(client, client.SDiff(keys...))
		return err
	})
	return
}

func (r *redisImpl) SDiffStore(dst string, keys ...string) (count int, err error) {
	if len(keys) == 0 {
		return
	}

	err = r.do("SDIFFSTORE", func(client driver.Client) error {
		count, err = mustBeInt(client, client.SDiffStore(dst, keys...))
		return err
	})
	return
}

func (r *redisImpl) SInter(keys ...string) (members []BulkString, err error) {
	if len(keys) == 0 {
		return
	}

	err = r.do("SINTER", func(client driver.Client) error {
		members, err = mustBeBulkStrings(client, client.SInter(keys...))
		return err
	})
	return
}

func (r *redisImpl) SInterStore(dst string, keys ...string) (count int, err error) {
	if len(keys) == 0 {
		return
	}

	err = r.do("SINTERSTORE", func(client driver.Client) error {
		count, err = mustBeInt(client, client.SInterStore(dst, keys...))
		return err
	})
	return
}

func (r *redisImpl) SIsMember(key string, member string) (exists bool, err error) {
	err = r.do("SISMEMBER", func(client driver.Client) error {
		exists, err = mustBeBool(client, client.SIsMember(key, member))
		return err
	})
	return
}

func (r *redisImpl) SMembers(key string) (members []BulkString, err error) {
	err = r.do("SMEMBERS", func(client driver.Client) error {
		members, err = mustBeBulkStrings(client, client.SMembers(key))
		return err
	})
	return
}

func (r *redisImpl) SMove(src string, dst string, member string) (moved bool, err error) {
	err = r.do("SMOVE", func(client driver.Client) error {
		moved, err = mustBeBool(client, client.SMove(src, dst, member))
		return err
	})
	return
}

func (r *redisImpl) SPop(key string) (member BulkString, err error) {
	err = r.do("SPOP", func(client driver.Client) error {
		member, err = mustBeBulkString(client, client.SPop(key))
		return err
	})
	return
}

func (r *redisImpl) SPopN(key string, count int) (members []BulkString, err error) {
	err = r.do("SPOP-N", func(client driver.Client) error {
		members, err = mustBeBulkStrings(client, client.SPopN(key, int64(count)))
		return err
	})
	return
}

func (r *redisImpl) SRandMember(key string) (member BulkString, err error) {
	err = r.do("SRANDMEMBER", func(client driver.Client) error {
		member, err = mustBeBulkString(client, client.SRandMember(key))
		return err
	})
	return
}

func (r *redisImpl) SRandMemberN(key string, count int) (members []BulkString, err error) {
	err = r.do("SRANDMEMBER-N", func(client driver.Client) error {
		members, err = mustBeBulkStrings(client, client.SRandMemberN(key, int64(count)))
		return err
	})
	return
}

func (r *redisImpl) SRem(key string, members ...string) (removed int, err error) {
	if len(members) == 0 {
		return
	}

	data := make([]interface{}, 0, len(members))

	for _, m := range members {
		data = append(data, m)
	}

	err = r.do("SREM", func(client driver.Client) error {
		removed, err = mustBeInt(client, client.SRem(key, data...))
		return err
	})
	return
}

func (r *redisImpl) SUnion(keys ...string) (members []BulkString, err error) {
	if len(keys) == 0 {
		return
	}

	err = r.do("SUNION", func(client driver.Client) error {
		members, err = mustBeBulkStrings(client, client.SUnion(keys...))
		return err
	})
	return
}

func (r *redisImpl) SUnionStore(dst string, keys ...string) (count int, err error) {
	if len(keys) == 0 {
		return
	}

	err = r.do("SUNIONSTORE", func(client driver.Client) error {
		count, err = mustBeInt(client, client.SUnionStore(dst, keys...))
		return err
	})
	return
}
