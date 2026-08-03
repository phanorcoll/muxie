package tmux

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandHomeDir(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("could not get home dir: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "expands tilde prefix",
			input:    "~/projects/muxie",
			expected: filepath.Join(home, "projects/muxie"),
		},
		{
			name:     "leaves absolute path unchanged",
			input:    "/usr/local/bin",
			expected: "/usr/local/bin",
		},
		{
			name:     "leaves relative path unchanged",
			input:    "projects/muxie",
			expected: "projects/muxie",
		},
		{
			name:     "empty string stays empty",
			input:    "",
			expected: "",
		},
		{
			name:     "tilde alone becomes home",
			input:    "~",
			expected: home,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandHomeDir(tt.input)
			if got != tt.expected {
				t.Errorf("expandHomeDir(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
