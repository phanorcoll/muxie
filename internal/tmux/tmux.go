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

// newWindowArgs builds the argument list for a tmux new-window command.
func newWindowArgs(sessionName, windowName, directory string) []string {
	return []string{"new-window", "-t", sessionName, "-n", windowName, "-c", directory}
}

// NewWindow creates a new tmux window in the specified session and directory.
// sessionName: the name of the tmux session.
// windowName: the name for the new window.
// directory: the working directory for the new window.
func NewWindow(sessionName, windowName, directory string) error {
	cmd := exec.Command("tmux", newWindowArgs(sessionName, windowName, directory)...)
	return cmd.Run()
}

// splitWindowArgs builds the argument list for a tmux split-window command.
// directory may be empty; when non-empty, -c directory is included.
func splitWindowArgs(sessionName, windowName, layout, directory string) []string {
	var layoutFlag string
	switch layout {
	case "horizontal":
		layoutFlag = "-h"
	case "vertical":
		layoutFlag = "-v"
	default:
		layoutFlag = "-h" // Default to horizontal split
	}
	args := []string{"split-window", "-t", fmt.Sprintf("%s:%s", sessionName, windowName), layoutFlag}
	if directory != "" {
		args = append(args, "-c", directory)
	}
	return args
}

// SplitWindow splits the current tmux window in the specified session.
// sessionName: the name of the tmux session.
// windowName: the name of the window to split.
// layout: the layout type ("horizontal" or "vertical").
// directory: optional working directory for the new pane (passed as -c).
func SplitWindow(sessionName, windowName, layout, directory string) error {
	cmd := exec.Command("tmux", splitWindowArgs(sessionName, windowName, layout, directory)...)
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

// WindowData represents information about a tmux window, including its panes.
type WindowData struct {
	Name   string
	Layout string
	Panes  []PaneData
}

// PaneData represents information about a single tmux pane.
type PaneData struct {
	Command   string
	Directory string
}

// GetSessionWindows retrieves all windows and their panes for a given tmux session.
func GetSessionWindows(sessionName string) ([]WindowData, error) {
	cmd := exec.Command("tmux", "list-windows", "-t", sessionName,
		"-F", "#{window_index}||#{window_name}||#{window_layout}")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("could not list windows for session '%s': %w", sessionName, err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var windows []WindowData

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "||", 3)
		if len(parts) < 2 {
			continue
		}
		windowIndex := parts[0]
		windowName := parts[1]
		windowLayout := ""
		if len(parts) > 2 {
			windowLayout = parts[2]
		}

		// Get panes for this window
		paneCmd := exec.Command("tmux", "list-panes", "-t",
			fmt.Sprintf("%s:%s", sessionName, windowIndex),
			"-F", "#{pane_current_command}||#{pane_current_path}")
		paneOutput, paneErr := paneCmd.Output()
		var panes []PaneData
		if paneErr == nil {
			paneLines := strings.Split(strings.TrimSpace(string(paneOutput)), "\n")
			for _, pl := range paneLines {
				if pl == "" {
					continue
				}
				pp := strings.SplitN(pl, "||", 2)
				cmd := pp[0]
				dir := ""
				if len(pp) > 1 {
					dir = pp[1]
				}
				panes = append(panes, PaneData{Command: cmd, Directory: dir})
			}
		}

		windows = append(windows, WindowData{Name: windowName, Layout: windowLayout, Panes: panes})
	}

	return windows, nil
}

// SelectLayout applies a layout string to a tmux window.
// The layout string can be a raw tmux layout (as returned by #{window_layout})
// or a preset name like even-horizontal, even-vertical, tiled, etc.
func SelectLayout(sessionName, windowName, layout string) error {
	cmd := exec.Command("tmux", "select-layout", "-t",
		fmt.Sprintf("%s:%s", sessionName, windowName), layout)
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
