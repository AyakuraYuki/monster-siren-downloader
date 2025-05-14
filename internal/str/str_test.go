package str

import "testing"

func TestReplaceDotSuffixRune(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{s: "file.txt", want: "file.txt"},
		{s: "no_dot", want: "no_dot"},
		{s: "dir.", want: "dir_"},
	}
	for _, tt := range tests {
		if got := ReplaceDotSuffixRune(tt.s); got != tt.want {
			t.Errorf("ReplaceDotSuffixRune(%q) = %q, want %q", tt.s, got, tt.want)
		}
	}
}
