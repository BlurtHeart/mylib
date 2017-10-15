package strutils

import (
	"errors"
	"sort"
)

var (
	InvalidSyntax   = errors.New("invalid syntax")
	ValueOutOfRange = errors.New("value out of range")
)

func Atoi(str string) (int, error) {
	if len(str) < 1 {
		return 0, InvalidSyntax
	}

	// int scope
	max := 1<<31 - 1
	min := -1 << 31

	a := 0
	flag := 1
	var i int

	if str[0] == '-' {
		flag = -1
		i = 1
	} else if str[0] == '+' {
		i = 1
	}
	for ; i < len(str); i++ {
		r := int(str[i] - '0')
		if r < 0 || r > 9 {
			return a, InvalidSyntax
		}
		if (max-r)/10 < a {
			return a, ValueOutOfRange
		}
		if (min+r)/10 > a {
			return a, ValueOutOfRange
		}
		a = a*10 + r*flag
	}

	return a, nil
}

type sortRunes []rune

func (s sortRunes) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortRunes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortRunes) Len() int {
	return len(s)
}

// return alphabet order of s
// "badc" -> "abcd"
func SortString(s string) string {
	r := []rune(s)
	sort.Sort(sortRunes(r))
	return string(r)
}
