package redis

// KeyAndValues 代表一系列 K/V 组合。
type KeyAndValues []KeyAndValue

// Map 将 kvs 转化成一个 map，方便使用。
func (kvs KeyAndValues) Map() map[string]string {
	if len(kvs) == 0 {
		return nil
	}

	m := make(map[string]string, len(kvs))

	for _, fv := range kvs {
		m[fv.Key] = fv.Value
	}

	return m
}

// KeyAndValue 代表一对 K/V 组合。
type KeyAndValue struct {
	Key   string
	Value string
}

// MakeKeyAndValue 可以方便的创建一个 KeyAndValue 实例。
func MakeKeyAndValue(key string, value string) KeyAndValue {
	return KeyAndValue{
		Key:   key,
		Value: value,
	}
}
