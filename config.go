package redis

import (
	"time"
)

const (
	// DefaultDialTimeout 代表默认的连接超时。
	DefaultDialTimeout = 5 * time.Second

	// DefaultReadTimeout 代表默认的读超时。
	DefaultReadTimeout = 3 * time.Second

	// DefaultWriteTimeout 代表默认的写超时。
	DefaultWriteTimeout = 3 * time.Second
)

// Config 代表一个 Redis 连接池工厂配置。
//
// 由于 Redis 有多种不兼容的连接池，这个类型由一系列子配置组成，如果同时设置了多个配置，生效的优先顺序如下：
//     - ClientConfig
//     - ClusterConfig
//     - FailoverConfig
//     - RingConfig TODO:
//     - SentinelConfig TODO:
//     - UniversalConfig TODO:
type Config struct {
	Client   *ClientConfig   `config:"client"`   // Client 是直连模式的配置。
	Cluster  *ClusterConfig  `config:"cluster"`  // Cluster 是集群模式的配置。
	Failover *FailoverConfig `config:"failover"` // Failover 是 failover client 的配置。
}

// ClientConfig 代表 Redis 直连模式的配置。
// TODO: 补充更多配置项。
type ClientConfig struct {
	Addr     string `config:"addr"`     // Addr 配置 Redis 服务地址。
	Password string `config:"password"` // Password 配置连接 Redis 的密码。
	DB       int    `config:"db"`       // DB 配置连接上 Redis 后默认选择的数据库。

	DialTimeout  time.Duration `config:"dail_timeout"`  // DialTimeout 配置连接超时，默认是 DefaultDialTimeout。
	ReadTimeout  time.Duration `config:"read_timeout"`  // ReadTimeout 配置读超时，默认是 DefaultReadTimeout。
	WriteTimeout time.Duration `config:"write_timeout"` // WriteTimeout 配置写超时，默认是 DefaultWriteTimeout。

	PoolSize int `config:"pool_size"` // PoolSize 配置连接池大小。
}

// ClusterConfig 代表 Redis cluster 配置。
// TODO: 补充更多配置项。
type ClusterConfig struct {
	Addrs    []string `config:"addrs"`    // Addrs 配置 Redis cluster 地址。
	Password string   `config:"password"` // Password 配置连接 Redis 的密码。

	DialTimeout  time.Duration `config:"dail_timeout"`  // DialTimeout 配置连接超时，默认是 DefaultDialTimeout。
	ReadTimeout  time.Duration `config:"read_timeout"`  // ReadTimeout 配置读超时，默认是 DefaultReadTimeout。
	WriteTimeout time.Duration `config:"write_timeout"` // WriteTimeout 配置写超时，默认是 DefaultWriteTimeout。

	PoolSize int `config:"pool_size"` // PoolSize 配置连接池大小。
}

// FailoverConfig 代表 Redis failover client 配置。
// TODO: 补充更多配置项。
type FailoverConfig struct {
	MasterName    string   `config:"master_name"`    // MasterName 代表 master 结点的名字。
	SentinelAddrs []string `config:"sentinel_addrs"` // SentinelAddrs 是哨兵地址。
	Password      string   `config:"password"`       // Password 配置连接 Redis 的密码。

	DialTimeout  time.Duration `config:"dail_timeout"`  // DialTimeout 配置连接超时，默认是 DefaultDialTimeout。
	ReadTimeout  time.Duration `config:"read_timeout"`  // ReadTimeout 配置读超时，默认是 DefaultReadTimeout。
	WriteTimeout time.Duration `config:"write_timeout"` // WriteTimeout 配置写超时，默认是 DefaultWriteTimeout。

	PoolSize int `config:"pool_size"` // PoolSize 配置连接池大小。
}
