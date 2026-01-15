// Package tmux provides utilities for interacting with and managing tmux sessions.
package tmux

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// SessionData represents information about a tmux session,
// including its name and the number of windows it contains.
type SessionData struct {
	Name          string // Name of the tmux session
	NumberWindows int    // Number of windows in the session
}

// GetSessionsList retrieves a list of all tmux sessions along with their window counts.
// Returns a slice of SessionData and an error if the command fails.
func GetSessionsList() ([]SessionData, error) {
	cmd := exec.Command("tmux", "list-sessions", "-F", "#S")
	output, err := cmd.Output()
	if err != nil {
		log.Println("Error listing tmux sessions:", err)
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var sessions []SessionData
	for _, name := range lines {
		if name != "" {
			winCmd := exec.Command("tmux", "list-windows", "-t", name, "-F", "#W")
			winOutput, winErr := winCmd.Output()
			numWindows := 0
			if winErr == nil {
				winLines := strings.SplitSeq(strings.TrimSpace(string(winOutput)), "\n")
				for w := range winLines {
					if w != "" {
						numWindows++
					}
				}
			} else {
				log.Println("Error listing windows for session", name, ":", winErr)
			}
			sessions = append(sessions, SessionData{Name: name, NumberWindows: numWindows})
		}
	}
	return sessions, nil
}

// RenameSession renames an existing tmux session from oldName to newName.
// Returns an error if the command fails.
func RenameSession(oldName, newName string) error {
	cmd := exec.Command("tmux", "rename-session", "-t", oldName, newName)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// CreateSession creates a new tmux session with the given name and starting directory.
// Switches to the new session after creation.
// Returns an error if the command fails.
func CreateSession(name string, dirname string) error {
	dirname = expandHomeDir(dirname)
	cmd := exec.Command("tmux", "new-session", "-d", "-s", name, "-c", dirname, "-n main")
	if err := cmd.Run(); err != nil {
		log.Println("Error creating session:", err)
		return err
	}
	switchCmd := exec.Command("tmux", "switch-client", "-t", name)
	if err := switchCmd.Run(); err != nil {
		log.Println("Error switching to new session:", err)
		return err
	}
	return nil
}

// KillSession kills the tmux session with the given name.
// Returns an error if the command fails.
func KillSession(name string) error {
	cmd := exec.Command("tmux", "kill-session", "-t", name)
	if err := cmd.Run(); err != nil {
		log.Println("Error killing session:", err)
		return err
	}
	return nil
}

// KillWindow kills the tmux window with the given session name and window index.
// Returns an error if the command fails.
func KillWindow(sessionName string, index int) error {
	cmd := exec.Command("tmux", "kill-window", "-t", fmt.Sprintf("%s:%d", sessionName, index))
	if err := cmd.Run(); err != nil {
		log.Println("Error killing session:", err)
		return err
	}
	return nil
}

// GetActiveSession returns the name of the currently active tmux session.
// Returns an empty string if there is no active session.
func GetActiveSession() (string, error) {
	cmd := exec.Command("tmux", "display-message", "-p", "#S")
	output, err := cmd.Output()
	if err != nil {
		log.Println("Error getting active session:", err)
		return "", err
	}
	activeSession := strings.TrimSpace(string(output))
	if activeSession == "" {
		return "", nil // No active session
	}
	return activeSession, nil
}

// NewWindow creates a new tmux window in the specified session and directory.
// sessionName: the name of the tmux session.
// windowName: the name for the new window.
// directory: the working directory for the new window.
func NewWindow(sessionName, windowName, directory string) error {
	cmd := exec.Command("tmux", "new-window", "-t", sessionName, "-n", windowName, "-c", directory)
	return cmd.Run()
}

// SplitWindow splits the current tmux window in the specified session.
// sessionName: the name of the tmux session.
// windowName: the name of the window to split.
// layout: the layout type ("horizontal" or "vertical").
func SplitWindow(sessionName, windowName, layout string) error {
	var layoutFlag string
	switch layout {
	case "horizontal":
		layoutFlag = "-h"
	case "vertical":
		layoutFlag = "-v"
	default:
		layoutFlag = "-h" // Default to horizontal split
	}
	cmd := exec.Command("tmux", "split-window", "-t", fmt.Sprintf("%s:%s", sessionName, windowName), fmt.Sprintf("%s", layoutFlag))
	return cmd.Run()
}

// GetPaneBaseIndex returns the value of the global tmux variable 'pane-base-index'
func GetPaneBaseIndex() (int, error) {
	cmd := exec.Command("tmux", "show", "-g", "pane-base-index")
	result, err := cmd.Output()
	if err != nil {
		return -1, fmt.Errorf("Failed to run 'tmux show -g pane-base-index': %s", err)
	}

	result_string := strings.TrimSpace(string(result))
	result_trimmed := strings.TrimPrefix(result_string, "pane-base-index ")

	i, err := strconv.Atoi(result_trimmed)
	if err != nil {
		return -1, fmt.Errorf("Failed to convert pane-base-index to integer %s", err)
	}

	return i, nil
}

// SendKeys sends a command to a specific tmux pane.
// sessionName: the name of the tmux session.
// windowName: the name of the window.
// paneIndex: the index of the pane (0-based).
// keys: the command or keys to send.
func SendKeys(sessionName, windowName string, paneIndex int, keys string) error {
	target := fmt.Sprintf("%s:%s.%d", sessionName, windowName, paneIndex)
	cmd := exec.Command("tmux", "send-keys", "-t", target)
	if keys != "" {
		cmd.Args = append(cmd.Args, keys, "C-m")
	}
	return cmd.Run()
}

// SwitchSession switches the tmux client to the specified session.
// sessionName: the name of the tmux session to switch to.
func SwitchSession(sessionName string) error {
	cmd := exec.Command("tmux", "switch-client", "-t", sessionName)
	return cmd.Run()
}

// expandHomeDir expands a leading ~ in a directory path to the user's home directory.
// dirname: the directory path to expand.
func expandHomeDir(dirname string) string {
	if dirname != "" && dirname[:1] == "~" {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, dirname[1:])
		}
	}
	return dirname
}
