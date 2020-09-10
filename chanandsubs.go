package redis

// ChanAndSubs 代表一系列 channel 和 subscriber 组合。
type ChanAndSubs []ChanAndSub

// Map 将 kvs 转化成一个 map，方便使用。
func (css ChanAndSubs) Map() map[string]int {
	if len(css) == 0 {
		return nil
	}

	m := make(map[string]int, len(css))

	for _, cs := range css {
		m[cs.Chan] = cs.Sub
	}

	return m
}

// ChanAndSub 代表一对 channel 和 subscriber。
type ChanAndSub struct {
	Chan string
	Sub  int
}

// MakeChanAndSub 可以方便的创建一个 ChanAndSub 实例。
func MakeChanAndSub(ch string, sub int) ChanAndSub {
	return ChanAndSub{
		Chan: ch,
		Sub:  sub,
	}
}
