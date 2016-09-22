package hunter

func AdvanceFilter(s string, suffixLen int) (bool, int) {
	// 出现的字符种类上限
	pass, n := OccurCharLimit(s, suffixLen)
	if !pass {
		return false, 0
	}

	return true, n
}

func OccurCharLimit(s string, suffixLen int) (bool, int) {
	d := make(map[rune]bool)

	body := s[:len(s)-suffixLen]
	for _, u := range body {
		d[u] = true
	}

	return len(d) <= Conf.AdvFilter.OccurChars, len(d)
}
