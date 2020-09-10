# go-redis：Redis 客户端 #

`go-redis` 封装了 Redis 客户端，并且提供了各种 Redis 命令的 Go 接口。

## 使用方法 ##

### 配置 `Redis` ###

如果使用 `go-runner` 启动服务，`go-redis` 会被自动创建，无需操心。所有配置放在配置文件的相关条目下。

注意，Redis 有多种启动模式，我们应该选用其中一种来填写配置：

* `[redis.client]`：代表普通直连模式，只在自测阶段应该设置，用于直连一台独立的 Redis 服务。
* `[redis.cluster]`：使用 Redis cluster，线上会使用 cluster 模式管理 Redis 集群。

以 `[redis.client]` 为例，配置内容如下：

```ini
[redis.client]
addr = "127.0.0.1:6379"
```

业务代码需要使用 Redis 时，直接使用 `New` 方法即可。

```go
import "github.com/altstory/go-redis"

func Foo(ctx context.Context, req *FooRequest) (res *FooResponse, err error) {
    // 这里省略各种参数检查……

    r := redis.New(ctx)
    v, err := r.Get("foo")

    if err != nil {
        // Redis 出错啦，记得写日志并报错。
        return
    }

    // 没取到 foo 的值，进行处理
    if v.IsNull() {
        // 处理……
    } else {
        // 拿到了 v 的值，通过 v.String() 得到里面的字符串内容。
    }

    // 其他代码省略……
}
```

### 在服务中使用多个 MySQL 连接 ###

在某些场景下，仅使用一个 Redis 并不足够，那么我们可以自行构建 `Factory` 来连接更多的 Redis 服务。

首先在配置文件中写一个新的 Redis 连接配置。

```ini
[redis_another.client]
addr = "127.0.0.1:6380"
```

然后实例化一个新的工厂。

```go
// anotherRedisFactory 的类型是 **Factory，是一个指针的指针。
var anotherRedisFactory = redis.Register("redis_another")
```

接着，使用这个全局变量 `anotherRedisFactory` 来创建新的 Redis client，在业务中使用。

```go
factory := *anotherRedisFactory
r := factory.New(ctx)

// 使用 redis client 进行各种操作……
```
