package redis

import (
	"context"
	"testing"

	"github.com/huandu/go-assert"
)

func TestHashMethods(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()
	f := factory(t)
	r := f.New(ctx)
	resetRedis(t, r)

	key := "hash-key"
	field1 := "field-1"
	value1 := "value-1"
	field2 := "field-2"
	value2 := "value-2"
	field3 := "field-3"
	value3 := "value-3"
	field4 := "field-4"

	isNew, err := r.HSet(key, field1, value2)
	a.NilError(err)
	a.Assert(isNew)

	isNew, err = r.HSet(key, field1, value1)
	a.NilError(err)
	a.Assert(!isNew)

	isNew, err = r.HSetNX(key, field1, value2)
	a.NilError(err)
	a.Assert(!isNew)

	v1, err := r.HGet(key, field1)
	a.NilError(err)
	a.Equal(v1.String(), value1)

	v4, err := r.HGet(key, field4)
	a.NilError(err)
	a.Assert(v4.IsNull())

	err = r.HMSet(key, MakeKeyAndValue(field2, value2), MakeKeyAndValue(field3, value3))
	a.NilError(err)

	values, err := r.HMGet(key, field1, field2, field3, field4)
	a.NilError(err)
	a.Equal(values, []BulkString{MakeBulkString(value1), MakeBulkString(value2), MakeBulkString(value3), Null()})

	exists, err := r.HExists(key, field1)
	a.NilError(err)
	a.Assert(exists)

	exists, err = r.HExists(key, field4)
	a.NilError(err)
	a.Assert(!exists)

	deleted, err := r.HDel(key, field1, field4)
	a.NilError(err)
	a.Equal(deleted, 1)

	exists, err = r.HExists(key, field1)
	a.NilError(err)
	a.Assert(!exists)
}
