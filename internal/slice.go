package internal

import (
	"strings"
)

// IsStringSliceContains --
func IsStringSliceContains(sl []string, s string) bool {
	for _, i := range sl {
		if i == s {
			return true
		}
	}

	return false
}

// TrimedSpaceStringSlice -- slice string by seperate and trim space
func TrimedSpaceStringSlice(s, sep string) []string {
	var sl []string

	for _, p := range strings.Split(s, sep) {
		if str := strings.TrimSpace(p); len(str) > 0 {
			sl = append(sl, strings.TrimSpace(p))
		}
	}

	return sl
}
