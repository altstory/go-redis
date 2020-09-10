package redis

// MemberAndScores 代表一系列 sorted set 的 member 和 score 组合。
type MemberAndScores []MemberAndScore

// Map 将 mss 转化成一个 map，方便使用。
func (mss MemberAndScores) Map() map[string]float64 {
	if len(mss) == 0 {
		return nil
	}

	m := make(map[string]float64, len(mss))

	for _, ms := range mss {
		m[ms.Member] = ms.Score
	}

	return m
}

// MemberAndScore 代表 sorted set 的 member 和 score。
type MemberAndScore struct {
	Member string
	Score  float64
}

// MakeMemberAndScore 可以方便的创建一个 MemberAndScore 实例。
func MakeMemberAndScore(member string, score float64) MemberAndScore {
	return MemberAndScore{
		Member: member,
		Score:  score,
	}
}
