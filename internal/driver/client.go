package driver

import (
	"github.com/go-redis/redis"
)

// Client 代表一个 Redis 客户端连接池。
// TODO: 暂时由于时间关系，直接暴露了 github.com/go-redis/redis 的接口，未来会改掉。
type Client interface {
	redis.Cmdable
	Close() error

	// 临时对外暴露这个接口，为了方便自定义命令。
	Process(cmder redis.Cmder) error

	// TODO: 还需要定义各种维护用的接口。
}
