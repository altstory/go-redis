package redis

// BulkString 代表 Redis 协议中的 Bulk String 类型，
// 详见 https://redis.io/topics/protocol#bulk-string-reply。
type BulkString struct {
	str *string
}

// MakeBulkString 将 str 包装成 BulkString。
func MakeBulkString(str string) BulkString {
	return BulkString{
		str: &str,
	}
}

// Null 返回一个 Null 值的 BulkString。
func Null() BulkString {
	return BulkString{}
}

// IsNull 判断 bs 是否为一个 Null 值。
func (bs BulkString) IsNull() bool {
	return bs.str == nil
}

// String 返回 bs 内部的字符串，如果 bs 为空，返回空字符串。
func (bs BulkString) String() string {
	if bs.IsNull() {
		return ""
	}

	return *bs.str
}
