package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis"

	"github.com/altstory/go-log"
	"github.com/altstory/go-redis/internal/driver"
	"github.com/altstory/go-runner"
)

var (
	defaultFactory = Register("redis")
)

// Factory 管理 Redis 连接池，并提供接口从连接池中取出可用的 Redis 连接。
type Factory struct {
	unavailable bool // 用来标记 Factory 是否完全不可用，方便 Register 能安全的工作。

	addrs  []string
	client driver.Client
	tested bool
}

// NewFactory 创建一个新的 Redis 连接池。
func NewFactory(config *Config) *Factory {
	var client driver.Client
	var addrs []string

	if config.Client != nil {
		addrs = []string{config.Client.Addr}
		client = newClientFromClientConfig(config.Client)
	} else if config.Cluster != nil {
		addrs = append(addrs, config.Cluster.Addrs...)
		client = newClientFromClusterConfig(config.Cluster)
	}

	return &Factory{
		addrs:  addrs,
		client: client,
	}
}

func newClientFromClientConfig(c *ClientConfig) driver.Client {
	dialTimeout := c.DialTimeout
	readTimeout := c.ReadTimeout
	writeTimeout := c.WriteTimeout

	if dialTimeout == 0 {
		dialTimeout = DefaultDialTimeout
	}

	if readTimeout == 0 {
		readTimeout = DefaultReadTimeout
	}

	if writeTimeout == 0 {
		writeTimeout = DefaultWriteTimeout
	}

	return redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       c.DB,

		DialTimeout:  dialTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,

		PoolSize: c.PoolSize,
	})
}

func newClientFromClusterConfig(c *ClusterConfig) driver.Client {
	dialTimeout := c.DialTimeout
	readTimeout := c.ReadTimeout
	writeTimeout := c.WriteTimeout

	if dialTimeout == 0 {
		dialTimeout = DefaultDialTimeout
	}

	if readTimeout == 0 {
		readTimeout = DefaultReadTimeout
	}

	if writeTimeout == 0 {
		writeTimeout = DefaultWriteTimeout
	}

	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    c.Addrs,
		Password: c.Password,

		DialTimeout:  dialTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,

		PoolSize: c.PoolSize,
	})
}

func newClientFromFailoverConfig(c *FailoverConfig) driver.Client {
	dialTimeout := c.DialTimeout
	readTimeout := c.ReadTimeout
	writeTimeout := c.WriteTimeout

	if dialTimeout == 0 {
		dialTimeout = DefaultDialTimeout
	}

	if readTimeout == 0 {
		readTimeout = DefaultReadTimeout
	}

	if writeTimeout == 0 {
		writeTimeout = DefaultWriteTimeout
	}

	return redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    c.MasterName,
		SentinelAddrs: c.SentinelAddrs,
		Password:      c.Password,

		DialTimeout:  dialTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,

		PoolSize: c.PoolSize,
	})
}

// Conn 连接 Redis 服务器并测试其可用性。
func (f *Factory) Conn(ctx context.Context) error {
	if f.unavailable {
		return errors.New("go-redis: factory is not initialized")
	}

	if f.client == nil {
		return errors.New("go-redis: factory is not initialized")
	}

	if err := f.client.Ping().Err(); err != nil {
		return errors.New("go-redis: fail to connect Redis")
	}

	f.tested = true
	return nil
}

// New 返回连接池中的一个连接。
func (f *Factory) New(ctx context.Context) Redis {
	if f.unavailable {
		return nil
	}

	if !f.tested {
		log.Errorf(ctx, "addrs=%v||go-redis: Redis factory is not connected (forgot to call `f.Conn`?)", f.addrs)
		return nil
	}

	return newRedis(ctx, f.client)
}

// Close 关闭连接池。
func (f *Factory) Close() error {
	if f.unavailable || f.client == nil {
		return nil
	}

	return f.client.Close()
}

// Register 将配置文件里 [section] 部分的配置用于初始化 Redis。
// 需要注意，Register 函数依赖于 runner 的启动流程，
// 在 AddClient 周期结束前，返回的 Factory 并不可用。
func Register(section string) **Factory {
	factory := &Factory{
		unavailable: true,
	}

	runner.AddClient(section, func(ctx context.Context, config *Config) error {
		if config == nil {
			return fmt.Errorf("go-redis: missing Redis config `[%v]`", section)
		}

		if config.Client == nil && config.Cluster == nil && config.Failover == nil {
			return fmt.Errorf("go-redis: fail to init Redis as there is no valid config in `[%v]`", section)
		}

		f := NewFactory(config)

		if err := f.Conn(ctx); err != nil {
			if config.Client != nil {
				log.Errorf(ctx, "err=%v||addr=%v||section=%v||go-redis: fail to init Redis in pooled client mode", err, config.Client.Addr, section)
			} else if config.Cluster != nil {
				log.Errorf(ctx, "err=%v||addr=%v||section=%v||go-redis: fail to init Redis in cluster mode", err, config.Cluster.Addrs, section)
			} else if config.Failover != nil {
				log.Errorf(ctx, "err=%v||master_name=%v||sentinel_addrs=%v||section=%v||go-redis: fail to init Redis in failover mode", err, config.Failover.MasterName, config.Failover.SentinelAddrs, section)
			}

			return err
		}

		if config.Client != nil {
			log.Tracef(ctx, "addr=%v||section=%v||go-redis: redis is connected", config.Client.Addr, section)
		} else if config.Cluster != nil {
			log.Tracef(ctx, "addr=%v||section=%v||go-redis: redis is connected", config.Cluster.Addrs, section)
		} else if config.Failover != nil {
			log.Tracef(ctx, "master_name=%v||sentinel_addrs=%v||section=%v||go-redis: redis is connected", config.Failover.MasterName, config.Failover.SentinelAddrs, section)
		}

		initMetrics()
		factory = f
		return nil
	})

	return &factory
}
