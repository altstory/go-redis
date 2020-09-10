package redis

import (
	"context"
	"testing"
)

func TestFactory(t *testing.T) {
	ctx := context.Background()
	f := factory(t)
	r := f.New(ctx)

	if r == nil {
		t.Fatalf("fail to create Redis client.")
	}

	resetRedis(t, r)

	const setKey = "foo"
	const setValue = "bar"

	if _, err := r.Set(setKey, setValue); err != nil {
		t.Fatalf("fail to SET.")
	}

	if v, err := r.Get(setKey); err != nil {
		t.Fatalf("fail to GET.")
	} else if v.String() != setValue {
		t.Fatalf("invalid GET result. [expected:%v] [actual:%v]", setValue, v)
	}
}
