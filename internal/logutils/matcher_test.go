package logutils_test

import (
	"testing"

	"github.com/Drathveloper/uslogs/internal/logutils"
)

func TestMasker_Mask(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		patterns []logutils.MaskPattern
		want     string
	}{
		{
			name:     "0 defined patterns shouldn't mask",
			input:    "foo bar baz",
			patterns: []logutils.MaskPattern{},
			want:     "foo bar baz",
		},
		{
			name:     "1 defined pattern should mask",
			input:    "foo:[bar baz]",
			patterns: []logutils.MaskPattern{logutils.NewMaskPattern("foo:[", '*', ']')},
			want:     "foo:[*******]",
		},
		{
			name:     "2 defined patterns should mask",
			input:    "foo:[bar baz] qux:[quux corge]",
			patterns: []logutils.MaskPattern{logutils.NewMaskPattern("foo:[", '*', ']'), logutils.NewMaskPattern("qux:[", '*', ']')},
			want:     "foo:[*******] qux:[**********]",
		},
		{
			name:     "pattern with no match delimiter shouldn't mask",
			input:    "foo:[bar baz]",
			patterns: []logutils.MaskPattern{logutils.NewMaskPattern("foo:", '*', '|')},
			want:     "foo:[bar baz]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dict := make([]string, 0, len(tt.patterns))
			for _, pattern := range tt.patterns {
				dict = append(dict, pattern.Start)
			}

			masker := logutils.NewMasker(dict...)

			got := masker.Mask([]byte(tt.input), tt.patterns)

			if string(got) != tt.want {
				t.Errorf("Masker.Mask() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
