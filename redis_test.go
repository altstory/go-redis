package redis

import (
	"context"
	"testing"
)

const testAddr = "127.0.0.1:6379"

func factory(t *testing.T) *Factory {
	f := NewFactory(&Config{
		Client: &ClientConfig{
			Addr: testAddr,
		},
	})

	ctx := context.Background()

	if err := f.Conn(ctx); err != nil {
		t.Skipf("fail to connect Redis server. [err:%v]", err)
	}

	return f
}

func resetRedis(t *testing.T, r Redis) {
	err := r.FlushAll()

	if err != nil {
		t.Fatalf("fail to flush redis. [err:%v]", err)
	}
}
