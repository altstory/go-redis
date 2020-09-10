package redis

import "github.com/altstory/go-redis/internal/driver"

// Connection 代表 Redis 跟连接相关的接口，详见 https://redis.io/commands#connection。
//
// 注意：由于 `SELECT` 和 `SWAPDB` 在 Redis Cluster 模式下已经没有用，
// 我们应该始终会使用 Redis Cluster 或者其他 Redis proxy 来实现高可用，
// 可以预计这两个接口应该永远不会用到，因此无需定义。
//
// 另外，`AUTH` 接口不做封装，密码验证的工作已经在框架层完成。
type Connection interface {
	Echo(msg string) (echo BulkString, err error)
	Ping() (err error)
}

func (r *redisImpl) Echo(msg string) (echo BulkString, err error) {
	err = r.do("ECHO", func(client driver.Client) error {
		echo, err = mustBeBulkString(client, client.Echo(msg))
		return err
	})
	return
}

func (r *redisImpl) Ping() (err error) {
	err = r.do("PING", func(client driver.Client) error {
		_, err = mustBeStatus(client, client.Ping())
		return err
	})
	return err
}
