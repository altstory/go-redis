package redis

import "github.com/altstory/go-redis/internal/driver"

// Server 代表 Redis 跟 server 相关的接口，详见 https://redis.io/commands#server。
type Server interface {
	// TODO:

	FlushAll(options ...FlushOption) (err error)
}

func (r *redisImpl) FlushAll(options ...FlushOption) (err error) {
	async := false

	for _, opt := range options {
		switch opt {
		case flushOptionAsync:
			async = true
		}
	}

	if async {
		err = r.do("FLUSHALL ASYNC", func(client driver.Client) error {
			_, err = mustBeStatus(client, client.FlushAllAsync())
			return err
		})
	} else {
		err = r.do("FLUSHALL", func(client driver.Client) error {
			_, err = mustBeStatus(client, client.FlushAll())
			return err
		})
	}

	return
}
