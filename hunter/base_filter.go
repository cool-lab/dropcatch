package hunter

import (
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/uber-go/zap"
)

func BaseFilter(s string) (bool, int) {
	// 过滤后缀类型，例如是否是.com域名
	if !SuffixTypePass(s) {
		return false, 0
	}

	// 过滤域名长度
	i := strings.Index(s, ".")
	if i == -1 {
		return false, 0
	}
	suffixLen := len(s) - i
	if !LenPass(s, suffixLen) {
		return false, 0
	}

	// '-'分隔符过滤
	if !DelimPass(s) {
		return false, 0
	}

	// 纯数字 OR 纯字母 OR 数字、字母都可以有
	if !CharsPass(s, suffixLen) {
		return false, 0
	}

	return true, suffixLen
}

func SuffixTypePass(s string) bool {
	return strings.HasSuffix(s, Conf.BaseFilter.SuffixType)
}

func LenPass(s string, suffixLen int) bool {
	return len(s)-suffixLen <= Conf.BaseFilter.MaxLen
}

func DelimPass(s string) bool {
	if Conf.BaseFilter.ExcludeDelim {
		// 不能包含 '-'
		return !strings.Contains(s, "-")
	}

	// 不限制
	return true
}

func CharsPass(s string, suffixLen int) bool {
	// 域名要是纯数字或纯字母
	body := s[:len(s)-suffixLen]

	switch Conf.BaseFilter.CharType {
	case 1: // 纯字母
		return govalidator.IsAlpha(body)
	case 2: // 纯数字
		return govalidator.IsNumeric(body)
	case 3:
		return govalidator.IsAlpha(body) || govalidator.IsNumeric(body)
	case 4:
		return true
	default:
		Logger.Fatal("invalid chartype", zap.Int("chartype", Conf.BaseFilter.CharType))
		return false
	}
}
