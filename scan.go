package redis

// Scan 代表所有跟扫描键值相关的接口，详见 https://redis.io/commands/scan。
type Scan interface {
	// TODO:
	// Scan(cursor int64, options ...ScanOption) (nextCursor int64, keys []string, err error)
	// SScan(key string, cursor int64, options ...ScanOption) (nextCursor int64, values []string, err error)
	// HScan(key string, cursor int64, options ...ScanOption) (nextCursor int64, fieldAndValues KeyAndValues, err error)
	// ZScan(key string, cursor int64, options ...ScanOption) (nextCursor int64, mss MemberAndScores, err error)
}
