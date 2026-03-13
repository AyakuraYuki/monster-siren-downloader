package filenames

import (
	"strings"
	"unicode/utf8"

	"github.com/flytam/filenamify"
)

const replacement = "_"

func SongName(name string) string {
	name, _ = filenamify.Filenamify(name, filenamify.Options{
		Replacement: replacement,
	})
	return strings.TrimSpace(name)
}

func AlbumName(name string) string {
	name, _ = filenamify.Filenamify(name, filenamify.Options{
		Replacement: replacement,
	})
	name = strings.TrimSpace(name)
	name = ReplaceDotSuffix(name)
	return strings.TrimSpace(name)
}

func ReplaceDotSuffix(s string) string {
	if s == "" {
		return s
	}
	lastRune, size := utf8.DecodeLastRuneInString(s)
	if lastRune == '.' {
		return s[:len(s)-size] + "_"
	}
	return s
}
