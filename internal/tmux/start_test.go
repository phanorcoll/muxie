package tmux

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePaneDirectory(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("could not get home dir: %v", err)
	}

	tests := []struct {
		name             string
		sessionDirectory string
		windowDirectory  string
		paneDirectory    string
		expected         string
	}{
		{
			name:             "uses pane directory when set",
			sessionDirectory: "/home/user/projects",
			windowDirectory:  "",
			paneDirectory:    "/home/user/other",
			expected:         "/home/user/other",
		},
		{
			name:             "falls back to window directory when pane not set",
			sessionDirectory: "/home/user/projects",
			windowDirectory:  "/home/user/window-special",
			paneDirectory:    "",
			expected:         "/home/user/window-special",
		},
		{
			name:             "falls back to session directory when neither set",
			sessionDirectory: "/home/user/projects",
			windowDirectory:  "",
			paneDirectory:    "",
			expected:         "/home/user/projects",
		},
		{
			name:             "pane directory wins over window directory",
			sessionDirectory: "/home/user/projects",
			windowDirectory:  "/home/user/window",
			paneDirectory:    "/home/user/pane",
			expected:         "/home/user/pane",
		},
		{
			name:             "expands tilde in pane directory",
			sessionDirectory: "/home/user/projects",
			windowDirectory:  "",
			paneDirectory:    "~/my-pane",
			expected:         filepath.Join(home, "my-pane"),
		},
		{
			name:             "expands tilde in window directory",
			sessionDirectory: "/home/user/projects",
			windowDirectory:  "~/my-window",
			paneDirectory:    "",
			expected:         filepath.Join(home, "my-window"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolvePaneDirectory(tt.sessionDirectory, tt.windowDirectory, tt.paneDirectory)
			if got != tt.expected {
				t.Errorf("resolvePaneDirectory(%q, %q, %q) = %q, want %q",
					tt.sessionDirectory, tt.windowDirectory, tt.paneDirectory, got, tt.expected)
			}
		})
	}
}

func TestSplitWindowCommandIncludesDirectory(t *testing.T) {
	tests := []struct {
		name          string
		sessionName   string
		windowName    string
		layout        string
		directory     string
		wantDirFlag   bool
	}{
		{
			name:        "includes -c when directory is provided",
			sessionName: "mysession",
			windowName:  "editor",
			layout:      "vertical",
			directory:   "/home/user/projects",
			wantDirFlag: true,
		},
		{
			name:        "omits -c when directory is empty",
			sessionName: "mysession",
			windowName:  "editor",
			layout:      "horizontal",
			directory:   "",
			wantDirFlag: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := splitWindowArgs(tt.sessionName, tt.windowName, tt.layout, tt.directory)

			// Verify -t flag is present
			foundTarget := false
			foundDir := false
			foundLayout := false

			for i, arg := range args {
				switch arg {
				case "-t":
					if i+1 < len(args) && args[i+1] == tt.sessionName+":"+tt.windowName {
						foundTarget = true
					}
				case "-c":
					if tt.wantDirFlag && i+1 < len(args) && args[i+1] == tt.directory {
						foundDir = true
					}
				case "-h", "-v":
					foundLayout = true
				}
			}

			if !foundTarget {
				t.Errorf("split-window args missing -t %s:%s in %v", tt.sessionName, tt.windowName, args)
			}
			if !foundLayout {
				t.Errorf("split-window args missing layout flag (-h or -v) in %v", args)
			}
			if tt.wantDirFlag && !foundDir {
				t.Errorf("split-window args missing -c %s in %v", tt.directory, args)
			}
			if !tt.wantDirFlag && foundDir {
				t.Errorf("split-window args should not contain -c when directory is empty, got %v", args)
			}
		})
	}
}

func TestNewWindowCommandIncludesDirectory(t *testing.T) {
	args := newWindowArgs("mysession", "editor", "/home/user/projects")

	if len(args) < 6 {
		t.Fatalf("expected at least 6 args, got %d: %v", len(args), args)
	}

	// new-window -t mysession -n editor -c /home/user/projects
	if args[0] != "new-window" {
		t.Errorf("expected 'new-window', got %q", args[0])
	}
	if args[1] != "-t" || args[2] != "mysession" {
		t.Errorf("expected '-t mysession', got %q %q", args[1], args[2])
	}
	if args[3] != "-n" || args[4] != "editor" {
		t.Errorf("expected '-n editor', got %q %q", args[3], args[4])
	}
	if args[5] != "-c" || args[6] != "/home/user/projects" {
		t.Errorf("expected '-c /home/user/projects', got %q %q", args[5], args[6])
	}
}
