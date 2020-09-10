package redis

import "strconv"

// ScoreRange 代表 sorted set 里面 score 的范围。
//
// 允许的取值为：
//     - 数字：-1、2.3
//     - 范围：(1
//     - 无穷：-inf、+inf
type ScoreRange string

// MakeScoreRange 将一个 score 变成 ScoreRange。
//
// 默认返回的范围是包含 score 这个值的，如果希望生成一个不包含 value 的 ScoreRange，应该调用
//     min := MakeScoreRange(1.234).Exclusive() // min == "(1.234"
func MakeScoreRange(score float64) ScoreRange {
	return ScoreRange(strconv.FormatFloat(score, 'f', -1, 64))
}

// Exclusive 将 sr 变成 exclusive 的值。
func (sr ScoreRange) Exclusive() ScoreRange {
	if sr == "" || sr[0] == '(' {
		return sr
	}

	return "(" + sr
}

// IsValid 判断 sr 是否合法。
func (sr ScoreRange) IsValid() bool {
	if sr == "" {
		return false
	}

	if c := sr[0]; c == '(' {
		sr = sr[1:]
	}

	_, err := strconv.ParseFloat(string(sr), 64)
	return err == nil
}

// String 返回 sr 的原值。
func (sr ScoreRange) String() string {
	return string(sr)
}

// MemberRange 表示一个 sorted set 里面 member 的范围。
//
// 允许的取值为：
//     - 包含字符串：[member1
//     - 不包含字符串：(member2
//     - 边界：-、+
type MemberRange string

// MemberRange 的边界值。
const (
	MemberRangeStart MemberRange = "-"
	MemberRangeStop  MemberRange = "+"
)

// MakeMemberRange 将一个 member 变成 MemberRange。
//
// 默认返回的范围是包含 member 这个值的，如果希望生成一个不包含 member 的 MemberRange，应该调用
//     min := MakeScoreRange(1.234).Exclusive() // min == "(1.234"
func MakeMemberRange(member string) MemberRange {
	return MemberRange("[" + member)
}

// Exclusive 将 mr 变成 exclusive 的值。
func (mr MemberRange) Exclusive() MemberRange {
	if mr == "" || mr[0] == '(' || mr == "-" || mr == "+" {
		return mr
	}

	if mr[0] == '[' {
		mr = mr[1:]
	}

	return "(" + mr
}

// IsValid 检查 lex 的合法性并返回结果。
// 根据 Redis 文档，lex 必须由 +、-、[ 或 ( 开头。
func (mr MemberRange) IsValid() bool {
	if mr == "" {
		return false
	}

	if mr == "+" || mr == "-" {
		return true
	}

	if c := mr[0]; c != '[' && c != '(' {
		return false
	}

	return true
}

// String 返回 mr 的字符串值。
func (mr MemberRange) String() string {
	return string(mr)
}
