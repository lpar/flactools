package flactools

import (
	"testing"
)

func TestCleanName(t *testing.T) {
	tests := []struct {
		Input  string
		Output string
	}{
		{"Mike's Album", "Mikes_Album"},
	}
	for _, test := range tests {
		got := CleanName(test.Input)
		if got != test.Output {
			t.Errorf("expected %s got %s", test.Output, got)
		}
	}
}
