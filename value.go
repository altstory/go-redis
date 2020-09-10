package redis

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/go-redis/redis"
)

var (
	// ErrKeyNotExist 是当 key 找不到时返回的错误。
	ErrKeyNotExist = errors.New("go-redis: key does not exist")

	// ErrKeyHasNoExpiration 是当 key 没有设置超时时间、使用 TTL/PTTL 时返回的错误。
	ErrKeyHasNoExpiration = errors.New("go-redis: key exists but has no associated expire")

	// ErrNotImplemented 表示这个功能还未实现。
	ErrNotImplemented = errors.New("go0redis: command is not implemented")

	// ErrUnexpectedResponseType 表示一个不支持的 redis 应答格式，一般都是这个库的 bug。
	ErrUnexpectedResponseType = errors.New("go-redis: unexpected type of the response")

	// ErrUnsupportedValueType 表示 redis 返回值的类型不支持，一般都是使用了还没实现的功能导致的。
	ErrUnsupportedValueType = errors.New("go-redis: unsupported value type")
)

// MultiValue 存储所有合法的 Redis 值类型，包括 nil 和 error。
// MULTI 和 EVAL 用它来存储返回值，调用者一般需要根据上下文来判断其中实际存储的类型是什么
type MultiValue struct {
	data interface{}
}

// MakeMultiValue 将 v 包装成一个 MultiValue。
// 如果 v 已经是一个 MultiValue 了，会返回 v 本身。
func MakeMultiValue(v interface{}) (mv MultiValue) {
	if v == nil {
		return
	}

	switch data := v.(type) {
	case error:
		if fmv, ok := data.(*FutureMultiValue); ok {
			mv = fmv.MultiValue()
		} else {
			mv.data = data
		}
	case []byte:
		if data == nil {
			mv.data = Null()
		} else {
			mv.data = MakeBulkString(string(data))
		}
	case string:
		mv.data = data
	case BulkString:
		mv.data = data
	case time.Duration:
		mv.data = data
	case time.Time:
		mv.data = data
	case bool:
		mv.data = data
	case int:
		mv.data = int64(data)
	case int64:
		mv.data = data
	case float64:
		mv.data = data
	case redis.Z:
		mv.data = MemberAndScore{
			Member: fmt.Sprint(data.Member),
			Score:  data.Score,
		}
	case redis.ClusterSlot:
		// FIXME: 支持 CLUSTER 相关命令的返回值。
		panic(ErrNotImplemented)
	case *redis.CommandInfo:
		// FIXME: 支持 COMMAND 相关命令的返回值。
		panic(ErrNotImplemented)
	case redis.GeoLocation:
		// FIXME: 支持 GEO 相关命令的返回值。
		panic(ErrNotImplemented)
	case *redis.GeoPos:
		// FIXME: 支持 GEO 相关命令的返回值。
		panic(ErrNotImplemented)

	case MultiValue:
		mv = data
	case KeyAndValue:
		mv.data = data
	case MemberAndScore:
		mv.data = data
	case ChanAndSub:
		mv.data = data

		// FIXME:
	// case *scanValues:
	// 	mv.data = data

	case []interface{}:
		mvs := make([]MultiValue, 0, len(data))

		for _, v := range data {
			mvs = append(mvs, MakeMultiValue(v))
		}

		mv.data = mvs
	case []error:
		mvs := make([]MultiValue, 0, len(data))

		for _, v := range data {
			mvs = append(mvs, MakeMultiValue(v))
		}

		mv.data = mvs
	case []string:
		mvs := make([]MultiValue, 0, len(data))

		for _, v := range data {
			mvs = append(mvs, MakeMultiValue(v))
		}

		mv.data = mvs
	case []bool:
		mvs := make([]MultiValue, 0, len(data))

		for _, v := range data {
			mvs = append(mvs, MakeMultiValue(v))
		}

		mv.data = mvs
	case map[string]string:
		mvs := make([]MultiValue, 0, len(data))

		for k, v := range data {
			mvs = append(mvs, MakeMultiValue(KeyAndValue{k, v}))
		}

		mv.data = mvs
	case map[string]int64:
		mvs := make([]MultiValue, 0, len(data))

		for k, v := range data {
			mvs = append(mvs, MakeMultiValue(ChanAndSub{
				Chan: k,
				Sub:  int(v),
			}))
		}

		mv.data = mvs
	case map[string]struct{}:
		mvs := make([]MultiValue, 0, len(data))

		for k := range data {
			mvs = append(mvs, MakeMultiValue(k))
		}

		mv.data = mvs
	case []redis.Z:
		mvs := make([]MultiValue, 0, len(data))

		for _, v := range data {
			mvs = append(mvs, MakeMultiValue(v))
		}

		mv.data = mvs
	case []redis.ClusterSlot:
		mvs := make([]MultiValue, 0, len(data))

		for _, v := range data {
			mvs = append(mvs, MakeMultiValue(v))
		}

		mv.data = mvs
	case map[string]*redis.CommandInfo:
		mvs := make([]MultiValue, 0, len(data))

		for _, v := range data {
			mvs = append(mvs, MakeMultiValue(v))
		}

		mv.data = mvs
	case []redis.GeoLocation:
		mvs := make([]MultiValue, 0, len(data))

		for _, v := range data {
			mvs = append(mvs, MakeMultiValue(v))
		}

		mv.data = mvs
	case []*redis.GeoPos:
		mvs := make([]MultiValue, 0, len(data))

		for _, v := range data {
			mvs = append(mvs, MakeMultiValue(v))
		}

		mv.data = mvs

	case []BulkString:
		mvs := make([]MultiValue, 0, len(data))

		for _, v := range data {
			mvs = append(mvs, MakeMultiValue(v))
		}

		mv.data = mvs
	case []MultiValue:
		mvs := make([]MultiValue, 0, len(data))

		for _, v := range data {
			mvs = append(mvs, v)
		}

		mv.data = mvs
	case KeyAndValues:
		mvs := make([]MultiValue, 0, len(data))

		for _, v := range data {
			mvs = append(mvs, MakeMultiValue(v))
		}

		mv.data = mvs
	case MemberAndScores:
		mvs := make([]MultiValue, 0, len(data))

		for _, v := range data {
			mvs = append(mvs, MakeMultiValue(v))
		}

		mv.data = mvs
	case ChanAndSubs:
		mvs := make([]MultiValue, 0, len(data))

		for _, v := range data {
			mvs = append(mvs, MakeMultiValue(v))
		}

		mv.data = mvs

		// FIXME:
	// case []scanValues:
	// 	mvs := make([]MultiValue, 0, len(data))

	// 	for _, v := range data {
	// 		mvs = append(mvs, MakeMultiValue(v))
	// 	}

	// 	mv.data = mvs

	default:
		val := reflect.ValueOf(v)

		switch val.Kind() {
		case reflect.String:
			return MakeMultiValue(val.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return MakeMultiValue(val.Int())
		case reflect.Float32, reflect.Float64:
			return MakeMultiValue(val.Float())
		case reflect.Bool:
			return MakeMultiValue(val.Bool())
		default:
			panic(ErrUnsupportedValueType)
		}
	}

	return
}

// IsNil 判断值是否为 nil。
func (mv MultiValue) IsNil() bool {
	return mv.data == nil
}

// IsErr 判断值是否为一个 error，如果是 error，所有其他的方法都会返回 ok 为 false。
func (mv MultiValue) IsErr() bool {
	_, ok := mv.data.(error)
	return ok
}

// Err 会在 MultiValue 存储的是 error 的时候返回对应的实例，否则返回 nil。
func (mv MultiValue) Err() error {
	if err, ok := mv.data.(error); ok {
		return err
	}

	return nil
}

// String 返回 mv 的字符串表示，用于输出日志。
func (mv MultiValue) String() string {
	return fmt.Sprintf("redis.MultiValue{data=%v}", mv.data)
}

// Int 返回整数，如果 MultiValue 存储的并不是 int，ok 会返回 false。
// Int 不区分 int64 与 int，对于 int64 数据来说，会截断成 int 值。
func (mv MultiValue) Int() (n int, ok bool) {
	n, ok = mv.data.(int)

	if !ok {
		if i, okok := mv.data.(int64); okok {
			n = int(i)
			ok = true
		}
	}

	return
}

// Int64 返回整数，如果 MultiValue 存储的并不是 int64，ok 会返回 false。
// Int64 不区分 int64 与 int，对于 int 数据来说，会转化成 int64 值。
func (mv MultiValue) Int64() (n int64, ok bool) {
	n, ok = mv.data.(int64)

	if !ok {
		if i, okok := mv.data.(int); okok {
			n = int64(i)
			ok = true
		}
	}

	return
}

// Float64 返回 float64，如果 MultiValue 存储的类型不是 float64，ok 为 false。
func (mv MultiValue) Float64() (n float64, ok bool) {
	n, ok = mv.data.(float64)
	return
}

// Bool 返回 bool 值，如果 MultiValue 存储的类型不是 bool 或 int/int64，ok 为 false。如果类型是 int/int64，非 0 代表 true。
func (mv MultiValue) Bool() (v bool, ok bool) {
	v, ok = mv.data.(bool)

	if ok {
		return
	}

	n, ok := mv.Int64()

	if !ok {
		return
	}

	v = n != 0
	return
}

// Status 返回状态码，如果 MultiValue 存储的类型不是 Status，ok 为 false。
func (mv MultiValue) Status() (s string, ok bool) {
	s, ok = mv.data.(string)
	return
}

// BulkString 返回 bulk string，如果 MultiValue 存储的类型不是 bulk string，ok 为 false。
func (mv MultiValue) BulkString() (v BulkString, ok bool) {
	if mv.data == nil {
		ok = true
		return
	}

	v, ok = mv.data.(BulkString)

	// 由于 github.com/go-redis/redis 在底层不一定能精确区分 status 和 string，
	// 例如 EVAL 返回里面被统一强制转成 string 了，因此这里只好做一个兼容。
	// 所幸一般来说我们不会被 status 和 string 困扰，这两种类型基本不可能混用，
	// 所以这个兼容不会影响到业务代码。
	if !ok {
		s, okok := mv.data.(string)

		if okok {
			ok = true
			v = MakeBulkString(s)
		}
	}

	return
}

// MultiValues 返回 []MultiValue，如果 MultiValue 存储的类型不是 []MultiValue，ok 为 false。
func (mv MultiValue) MultiValues() (mvs []MultiValue, ok bool) {
	if mv.data == nil {
		ok = true
		return
	}

	mvs, ok = mv.data.([]MultiValue)
	return
}

// Time 返回一个 time.Time，如果 MultiValue 存储的类型不是 time.Time，ok 为 false。
func (mv MultiValue) Time() (t time.Time, ok bool) {
	t, ok = mv.data.(time.Time)
	return
}

// Duration 返回一个 time.Duration，如果 MultiValue 存储的类型不是 time.Duration，ok 为 false。
func (mv MultiValue) Duration() (d time.Duration, ok bool) {
	d, ok = mv.data.(time.Duration)
	return
}

// KeyAndValue 返回一个 KeyAndValue，如果 MultiValue 存储的类型不是 KeyAndValue，ok 为 false。
func (mv MultiValue) KeyAndValue() (kv KeyAndValue, ok bool) {
	kv, ok = mv.data.(KeyAndValue)
	return
}

// MemberAndScore 返回一个 MemberAndScore，如果 MultiValue 存储的类型不是 MemberAndScore，ok 为 false。
func (mv MultiValue) MemberAndScore() (ms MemberAndScore, ok bool) {
	ms, ok = mv.data.(MemberAndScore)
	return
}

func isPipelined(cmdable redis.Cmdable) (ok bool) {
	_, ok = cmdable.(redis.Pipeliner)
	return
}

func parsePipelinedReply(cmders []redis.Cmder, err error) ([]MultiValue, error) {
	if err != nil && err != redis.Nil {
		return nil, err
	}

	// 如果错误当前 err 是 redis.Nil，那么需要看看是否后面有其他错误存在，这才能确保所有都没有错。
	if err == redis.Nil {
		err = nil

		for _, cmder := range cmders {
			if e := cmder.Err(); e != nil && e != redis.Nil {
				err = e
				break
			}
		}
	}

	mvs := make([]MultiValue, 0, len(cmders))

	for _, cmder := range cmders {
		mv, e := parseCmder(cmder)

		if e != nil && err != nil {
			err = e
			mv = MakeMultiValue(e)
		}

		mvs = append(mvs, mv)
	}

	return mvs, nil
}

// parseCmder 将任意的 cmder 解析成可以直接被使用的 MultiValue.
// 需要注意的是，如果 cmder.Err() != nil，这里面的 error 会被存在 MultiValue 中，err 依然还是 nil。
// 只有当 cmder 无法被解析（例如当前还不认识的 cmder）才会返回 err。
func parseCmder(cmder redis.Cmder) (mv MultiValue, err error) {
	err = cmder.Err()

	if err == redis.Nil {
		err = nil

		if _, ok := cmder.(*redis.StringCmd); ok {
			mv = MakeMultiValue(Null())
		}

		return
	}

	if err != nil {
		mv = MakeMultiValue(err)
		err = nil
		return
	}

	switch c := cmder.(type) {
	case *redis.Cmd:
		mv = MakeMultiValue(c.Val())
	case *redis.SliceCmd:
		mv = MakeMultiValue(c.Val())
	case *redis.StatusCmd:
		s := c.Val()
		mv = MakeMultiValue(s)
	case *redis.IntCmd:
		mv = MakeMultiValue(c.Val())
	case *redis.DurationCmd:
		mv = MakeMultiValue(c.Val())
	case *redis.TimeCmd:
		mv = MakeMultiValue(c.Val())
	case *redis.BoolCmd:
		mv = MakeMultiValue(c.Val())
	case *redis.StringCmd:
		mv = MakeMultiValue(MakeBulkString(c.Val()))
	case *redis.FloatCmd:
		mv = MakeMultiValue(c.Val())
	case *redis.StringSliceCmd:
		mv = MakeMultiValue(c.Val())
	case *redis.BoolSliceCmd:
		mv = MakeMultiValue(c.Val())
	case *redis.StringStringMapCmd:
		mv = MakeMultiValue(c.Val())
	case *redis.StringIntMapCmd:
		mv = MakeMultiValue(c.Val())
	case *redis.StringStructMapCmd:
		mv = MakeMultiValue(c.Val())
	case *redis.ZSliceCmd:
		mv = MakeMultiValue(c.Val())

		// FIXME:
	// case *redis.ScanCmd:
	// 	keys, cursor := c.Val()
	// 	mv = MakeMultiValue(&scanValues{
	// 		cursor: int64(cursor),
	// 		values: makeValues(keys...),
	// 	})
	case *redis.ClusterSlotsCmd, *redis.GeoLocationCmd, *redis.GeoPosCmd, *redis.CommandsInfoCmd:
		err = ErrNotImplemented
	default:
		err = ErrNotImplemented
	}

	return
}

func mustBeMultiValue(cmdable redis.Cmdable, cmder redis.Cmder) (mv MultiValue, err error) {
	if isPipelined(cmdable) {
		err = &FutureMultiValue{cmder}
		return
	}

	v, e := parseCmder(cmder)

	if e != nil {
		panic(e)
	}

	if err = v.Err(); err != nil {
		return
	}

	mv = v
	return
}

func mustBeBool(cmdable redis.Cmdable, cmder redis.Cmder) (v bool, err error) {
	if isPipelined(cmdable) {
		err = &FutureMultiValue{cmder}
		return
	}

	mv, e := parseCmder(cmder)

	if e != nil {
		panic(e)
	}

	if err = mv.Err(); err != nil {
		return
	}

	v, ok := mv.Bool()

	if !ok {
		panic(ErrUnexpectedResponseType)
	}

	return
}

func mustBeInt(cmdable redis.Cmdable, cmder redis.Cmder) (n int, err error) {
	if isPipelined(cmdable) {
		err = &FutureMultiValue{cmder}
		return
	}

	mv, e := parseCmder(cmder)

	if e != nil {
		panic(e)
	}

	if err = mv.Err(); err != nil {
		return
	}

	v, ok := mv.Int()

	if !ok {
		panic(ErrUnexpectedResponseType)
	}

	n = v
	return
}

func mustBeIntOrNil(cmdable redis.Cmdable, cmder redis.Cmder) (n int, exists bool, err error) {
	if isPipelined(cmdable) {
		err = &FutureMultiValue{cmder}
		return
	}

	mv, e := parseCmder(cmder)

	if e != nil {
		panic(e)
	}

	if err = mv.Err(); err != nil {
		return
	}

	if mv.IsNil() {
		return
	}

	v, ok := mv.Int()

	if !ok {
		panic(ErrUnexpectedResponseType)
	}

	n = v
	exists = true
	return
}

func mustBeInt64(cmdable redis.Cmdable, cmder redis.Cmder) (n int64, err error) {
	if isPipelined(cmdable) {
		err = &FutureMultiValue{cmder}
		return
	}

	mv, e := parseCmder(cmder)

	if e != nil {
		panic(e)
	}

	if err = mv.Err(); err != nil {
		return
	}

	v, ok := mv.Int64()

	if !ok {
		panic(ErrUnexpectedResponseType)
	}

	n = v
	return
}

func mustBeFloat64(cmdable redis.Cmdable, cmder redis.Cmder) (n float64, err error) {
	if isPipelined(cmdable) {
		err = &FutureMultiValue{cmder}
		return
	}

	mv, e := parseCmder(cmder)

	if e != nil {
		panic(e)
	}

	if err = mv.Err(); err != nil {
		return
	}

	v, ok := mv.Float64()

	if !ok {
		panic(ErrUnexpectedResponseType)
	}

	n = v
	return
}

func mustBeStatus(cmdable redis.Cmdable, cmder redis.Cmder) (s string, err error) {
	if isPipelined(cmdable) {
		err = &FutureMultiValue{cmder}
		return
	}

	mv, e := parseCmder(cmder)

	if e != nil {
		panic(e)
	}

	if err = mv.Err(); err != nil {
		return
	}

	v, ok := mv.Status()

	if !ok {
		panic(ErrUnexpectedResponseType)
	}

	s = v
	return
}

func mustBeBulkString(cmdable redis.Cmdable, cmder redis.Cmder) (s BulkString, err error) {
	if isPipelined(cmdable) {
		err = &FutureMultiValue{cmder}
		return
	}

	mv, e := parseCmder(cmder)

	if e != nil {
		panic(e)
	}

	if err = mv.Err(); err != nil {
		return
	}

	v, ok := mv.BulkString()

	if !ok {
		panic(ErrUnexpectedResponseType)
	}

	s = v
	return
}

func mustBeMultiValues(cmdable redis.Cmdable, cmder redis.Cmder) (mvs []MultiValue, err error) {
	if isPipelined(cmdable) {
		err = &FutureMultiValue{cmder}
		return
	}

	mv, e := parseCmder(cmder)

	if e != nil {
		panic(e)
	}

	if err = mv.Err(); err != nil {
		return
	}

	v, ok := mv.MultiValues()

	if !ok {
		panic(ErrUnexpectedResponseType)
	}

	mvs = v
	return
}

func mustBeTime(cmdable redis.Cmdable, cmder redis.Cmder) (t time.Time, err error) {
	if isPipelined(cmdable) {
		err = &FutureMultiValue{cmder}
		return
	}

	mv, e := parseCmder(cmder)

	if e != nil {
		panic(e)
	}

	if err = mv.Err(); err != nil {
		return
	}

	tm, ok := mv.Time()

	if !ok {
		panic(ErrUnexpectedResponseType)
	}

	t = tm
	return
}

func mustBeDuration(cmdable redis.Cmdable, cmder redis.Cmder) (d time.Duration, err error) {
	if isPipelined(cmdable) {
		err = &FutureMultiValue{cmder}
		return
	}

	mv, e := parseCmder(cmder)

	if e != nil {
		panic(e)
	}

	if err = mv.Err(); err != nil {
		return
	}

	duration, ok := mv.Duration()

	if !ok {
		panic(ErrUnexpectedResponseType)
	}

	d = duration
	return
}

func mustBeKeyAndValue(cmdable redis.Cmdable, cmder redis.Cmder) (kv KeyAndValue, err error) {
	if isPipelined(cmdable) {
		err = &FutureMultiValue{cmder}
		return
	}

	mv, e := parseCmder(cmder)

	if e != nil {
		panic(e)
	}

	if err = mv.Err(); err != nil {
		return
	}

	val, ok := mv.KeyAndValue()

	if !ok {
		panic(ErrUnexpectedResponseType)
	}

	kv = val
	return
}

func mustBeMemberAndScore(cmdable redis.Cmdable, cmder redis.Cmder) (ms MemberAndScore, err error) {
	if isPipelined(cmdable) {
		err = &FutureMultiValue{cmder}
		return
	}

	mv, e := parseCmder(cmder)

	if e != nil {
		panic(e)
	}

	if err = mv.Err(); err != nil {
		return
	}

	val, ok := mv.MemberAndScore()

	if !ok {
		panic(ErrUnexpectedResponseType)
	}

	ms = val
	return
}

func mustBeBulkStrings(cmdable redis.Cmdable, cmder redis.Cmder) (bss []BulkString, err error) {
	mvs, e := mustBeMultiValues(cmdable, cmder)

	if e != nil {
		err = e
		return
	}

	bss = make([]BulkString, 0, len(mvs))

	for _, mv := range mvs {
		v, ok := mv.BulkString()

		if !ok {
			panic(ErrUnexpectedResponseType)
		}

		bss = append(bss, v)
	}

	return
}

func mustBeKeyAndValues(cmdable redis.Cmdable, cmder redis.Cmder) (kvs []KeyAndValue, err error) {
	mvs, e := mustBeMultiValues(cmdable, cmder)

	if e != nil {
		err = e
		return
	}

	kvs = make([]KeyAndValue, 0, len(mvs))

	for _, mv := range mvs {
		v, ok := mv.KeyAndValue()

		if !ok {
			panic(ErrUnexpectedResponseType)
		}

		kvs = append(kvs, v)
	}

	return
}

func mustBeMemberAndScores(cmdable redis.Cmdable, cmder redis.Cmder) (mss MemberAndScores, err error) {
	mvs, e := mustBeMultiValues(cmdable, cmder)

	if e != nil {
		err = e
		return
	}

	mss = make([]MemberAndScore, 0, len(mvs))

	for _, mv := range mvs {
		v, ok := mv.MemberAndScore()

		if !ok {
			panic(ErrUnexpectedResponseType)
		}

		mss = append(mss, v)
	}

	return
}

// FutureMultiValue 表示一个还未获得结果的 MultiValue，
// 一般来说使用者不需要直接用这个结构，而是把它当做 error 来用，
// 使用 MakeMultiValue 来还原里面的值。
//
// FutureMultiValue 实现了 error 接口，
// transaction 或 pipeline 中的所有 redis 接口调用都会返回一个 FutureMultiValue error，
// 使用者可以将这个 error 存起来，等 transaction 或 pipeline 提交之后
// 再用 MakeMultiValue 转化成真正的 MultiValue。
//
// 例如：
//     var futureValue error
//     _, err := client.Multi(r redis.Redis) error {
//         r.Set("foo", "2")
//         r.Incr("foo")
//         _, futureValue = r.Get("foo")
//     })
//     /* 检查 err，这里略过 */
//
//     // 使用 MakeMultiValue 获得真正的 MultiValue。
//     val := redis.MakeMultiValue(futureValue)
//     v, ok := val.Value()
//     fmt.Println(v, ok) // Output: 3 true
type FutureMultiValue struct {
	cmder redis.Cmder
}

// Error 返回错误信息。
func (fmv *FutureMultiValue) Error() string {
	return "redis: result is pipelined"
}

// MultiValue 返回内部的 MultiValue 值，如果值为空，
// 会返回一个 IsNil 为 true 的 MultiValue。
func (fmv *FutureMultiValue) MultiValue() MultiValue {
	mv, _ := parseCmder(fmv.cmder)
	return mv
}
