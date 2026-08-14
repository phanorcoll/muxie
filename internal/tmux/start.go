// Package tmux provides utilities for interacting with and managing tmux sessions.
package tmux

import (
	"fmt"
	"github.com/phanorcoll/muxie/internal/config"
)

// resolvePaneDirectory resolves the effective directory for a pane using the cascade:
// pane directory > window directory > session directory.
// Tilde paths are expanded via expandHomeDir.
func resolvePaneDirectory(sessionDir, windowDir, paneDir string) string {
	resolved := sessionDir
	if windowDir != "" {
		resolved = expandHomeDir(windowDir)
	}
	if paneDir != "" {
		resolved = expandHomeDir(paneDir)
	}
	return resolved
}

// StartSession creates a new tmux session with the given sessionName and starting directory.
// It then creates the specified windows and panes, running the configured commands in each pane.
// Returns an error if any tmux operation fails.
func StartSession(sessionName, sessionDirectory string, windows []config.Window) error {
	if err := CreateSession(sessionName, sessionDirectory); err != nil {
		return fmt.Errorf("failed to create new session '%s': %w", sessionName, err)
	}
	sessionDirectory = expandHomeDir(sessionDirectory)

	// Get the base pane index, which we will use several times
	// while starting a predefined session
	basePaneIndex, err := GetPaneBaseIndex()
	if err != nil {
		return err
	}

	for _, w := range windows {
		windowDirectory := sessionDirectory
		if w.Directory != "" {
			windowDirectory = expandHomeDir(w.Directory)
		}

		if err := NewWindow(sessionName, w.Name, windowDirectory); err != nil {
			return fmt.Errorf("failed to create window '%s' in session '%s': %w", w.Name, sessionName, err)
		}

		for j, p := range w.Panes {
			paneDirectory := resolvePaneDirectory(sessionDirectory, w.Directory, p.Directory)

			if j > 0 {
				// Use the simple direction if available; raw layout strings are
				// applied later via select-layout.
				splitLayout := w.Layout
				if splitLayout != "horizontal" && splitLayout != "vertical" {
					splitLayout = "horizontal" // placeholder; real layout applied below
				}
				if err := SplitWindow(sessionName, w.Name, splitLayout, paneDirectory); err != nil {
					return fmt.Errorf("failed to split window '%s' in session '%s': %w", w.Name, sessionName, err)
				}
			}

			if err := SendKeys(sessionName, w.Name, basePaneIndex + j, fmt.Sprintf("cd %s && clear && %s", paneDirectory, p.Command)); err != nil {
				return fmt.Errorf("failed to send keys to pane %d in window '%s' of session '%s': %w", j, w.Name, sessionName, err)
			}
		}

		// Apply the raw layout string if present — this handles nested splits
		// that can't be expressed with a simple -h/-v flag.
		if w.Layout != "" && w.Layout != "horizontal" && w.Layout != "vertical" {
			if err := SelectLayout(sessionName, w.Name, w.Layout); err != nil {
				return fmt.Errorf("failed to apply layout to window '%s': %w", w.Name, err)
			}
		}
	}

	// After starting the predefined session, delete its starting window
	// as it is not part of the user's config file
	err = KillWindow(sessionName, basePaneIndex)
	if err != nil {
		return err
	}

	return nil
}
