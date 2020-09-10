package redis

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/altstory/go-log"
	"github.com/altstory/go-redis/internal/driver"
)

const (
	logTagRedis = "_redis"
)

// Redis 代表一个 Redis 连接。
type Redis interface {
	Cluster
	Connection
	Generic
	GEO
	Hashes
	HyperLogLog
	Lists
	PubSub
	Scan
	Scripting
	Server
	Sets
	SortedSets
	Streams
	Strings
	Transactions
}

type redisImpl struct {
	// TODO: 真正实现这个接口
	Redis

	ctx    context.Context
	client driver.Client
}

var _ Redis = new(redisImpl)

// New 通过默认的连接池创建一个 Redis 连接。
func New(ctx context.Context) Redis {
	factory := *defaultFactory

	if factory == nil {
		log.Errorf(ctx, "go-redis: default factory is not initialized (forgot to set Redis config?)")
		return nil
	}

	return factory.New(ctx)
}

func newRedis(ctx context.Context, client driver.Client) *redisImpl {
	return &redisImpl{
		ctx:    ctx,
		client: client,
	}
}

func (r *redisImpl) do(cmd string, fn func(client driver.Client) error) (err error) {
	ctx := r.ctx
	ctx = log.WithTag(ctx, logTagRedis)

	if err = ctx.Err(); err != nil {
		log.Infof(ctx, "err=%v||cmd=%v||go-redis: context timeout", err, cmd)
		return
	}

	now := time.Now()

	defer func() {
		dur := time.Now().Sub(now)
		dur = dur.Round(time.Millisecond)
		proctime := dur.Seconds()

		if r := recover(); r != nil {
			log.Errorf(ctx, "err=%v||cmd=%v||proctime=%v||go-redis: caught a panic with call stack\n%v", r, cmd, proctime, string(debug.Stack()))
			err = fmt.Errorf("go-redis: caught a panic in `%v`", cmd)
		}

		if err == nil {
			log.Tracef(ctx, "cmd=%v||proctime=%.6f||go-redis: success", cmd, proctime)
		} else {
			log.Errorf(ctx, "err=%v||cmd=%v||proctime=%.6f||go-redis: failed", err, cmd, proctime)
		}

		statsForCall(ctx, err)
	}()

	err = fn(r.client)
	return
}
