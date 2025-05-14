package str

import "unicode/utf8"

func ReplaceDotSuffixRune(s string) string {
	if s == "" {
		return s
	}
	lastRune, size := utf8.DecodeLastRuneInString(s)
	if lastRune == '.' {
		return s[:len(s)-size] + "_"
	}
	return s
}
