// Package tui contains terminal user interface commands and related functionality.
package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/phanorcoll/muxie/internal/config"
	"github.com/phanorcoll/muxie/internal/tmux"
)

// Update is part of the Bubble Tea Model interface. It handles incoming messages
// and updates the model's state accordingly. Currently, it returns the model unchanged.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case m.showInput:
			switch {
			case key.Matches(msg, m.keys.Enter):
				if m.statusData.action == "a" && m.showInput {
					newName := m.sessionInput.Value()
					if newName == "" || newName == m.activeSession {
						m.sessionInput.Blur()
						m.sessionInput.Reset()
						m.showInput = false
						m.statusData.action = ""
						m.statusData.actionTitle = ""
						m.statusData.icon = ""
						m.statusData.color = "#FFF7DB"
						return m, nil
					}
					err := tmux.CreateSession(newName, "")
					if err != nil {
						m.logger.Printf("Error creating session %s: %v", newName, err)
					}
					return m, tea.Quit
				}
				if m.statusData.action == "d" && m.showInput {
					option := m.sessionInput.Value()
					if option == "y" {
						si := m.sessionList.SelectedItem().(session)
						err := tmux.KillSession(si.sessionName)
						if err != nil {
							m.logger.Printf("Error killing session %s: %v", si.sessionName, err)
						}
						if !si.isFromConfig {
							index := m.sessionList.GlobalIndex()
							m.sessionList.RemoveItem(index)
						}
					}
					m.sessionInput.Blur()
					m.sessionInput.Reset()
					m.showInput = false
					m.statusData.action = ""
					m.statusData.actionTitle = ""
					m.statusData.icon = ""
					m.statusData.color = "#FFF7DB"
					return m, getSessionsCmd(m.config)
				}
				if m.statusData.action == "r" && m.showInput {
					si := m.sessionList.SelectedItem().(session)
					newName := m.sessionInput.Value()
					if newName != "" {
						index := m.sessionList.GlobalIndex()
						m.sessionList.SetItem(index, session{
							sessionName: newName,
						})
						err := tmux.RenameSession(si.sessionName, newName)
						if err != nil {
							m.logger.Printf("Error renaming session to %s: %v", newName, err)
						}
					}
					m.sessionInput.Blur()
					m.sessionInput.Reset()
					m.showInput = false
					m.statusData.action = ""
					m.statusData.actionTitle = ""
					m.statusData.icon = ""
					m.statusData.color = "#FFF7DB"
					return m, getSessionsCmd(m.config)
				}

			case key.Matches(msg, m.keys.Escape):
				m.sessionInput.Blur()
				m.sessionInput.Reset()
				m.showInput = false
				m.statusData.action = ""
				m.statusData.actionTitle = ""
				m.statusData.icon = ""
				m.statusData.color = "#FFF7DB"
				return m, nil
			}
		default:
			if m.sessionList.FilterState() == list.Filtering {
				break
			}
			switch {
			case key.Matches(msg, m.keys.Add):
				m.showInput = true
				m.sessionInput.Placeholder = "type"
				m.sessionInput.Focus()
				m.statusData.actionTitle = "Name new session"
				m.statusData.action = "a"
				m.statusData.icon = ""
				m.statusData.color = "#37A1C5"
				return m, nil
			case key.Matches(msg, m.keys.Rename):
				si := m.sessionList.SelectedItem().(session)
				if si.isFromConfig {
					statusCmd := m.sessionList.NewStatusMessage(errorStyle("󰗼 rename in config file"))
					return m, statusCmd
				}
				m.showInput = true
				m.sessionInput.Placeholder = "type"
				m.sessionInput.Focus()
				m.statusData.actionTitle = fmt.Sprintf("Rename %s", si.sessionName)
				m.statusData.action = "r"
				m.statusData.icon = "󰑕"
				m.statusData.color = "#BAC537"
				return m, nil
			case key.Matches(msg, m.keys.Kill):
				si := m.sessionList.SelectedItem().(session)
				if !si.isRunning {
					statusCmd := m.sessionList.NewStatusMessage(errorStyle("󰗼 not running"))
					return m, statusCmd
				}
				m.showInput = true
				m.sessionInput.Placeholder = "y/n"
				m.sessionInput.Focus()
				m.statusData.action = "d"
				m.statusData.actionTitle = fmt.Sprintf("kill %s ?", si.sessionName)
				m.statusData.icon = "󰗨"
				m.statusData.color = "#C53770"
				return m, nil
			case key.Matches(msg, m.keys.Save):
				si := m.sessionList.SelectedItem().(session)
				if si.isFromConfig {
					statusCmd := m.sessionList.NewStatusMessage(errorStyle("󰗼 already in config"))
					return m, statusCmd
				}
				if !si.isRunning {
					statusCmd := m.sessionList.NewStatusMessage(errorStyle("󰗼 not running"))
					return m, statusCmd
				}
				windows, err := tmux.GetSessionWindows(si.sessionName)
				if err != nil {
					m.logger.Printf("Error getting windows for session %s: %v", si.sessionName, err)
					statusCmd := m.sessionList.NewStatusMessage(errorStyle("󰗼 failed to read session"))
					return m, statusCmd
				}
				var cfgWindows []config.Window
				for _, w := range windows {
					var panes []config.Pane
					for _, p := range w.Panes {
						panes = append(panes, config.Pane{
							Command:   p.Command,
							Directory: p.Directory,
						})
					}
					cfgWindows = append(cfgWindows, config.Window{
						Name:   w.Name,
						Layout: w.Layout,
						Panes:  panes,
					})
				}
				m.config.Sessions = append(m.config.Sessions, config.Session{
					Name:    si.sessionName,
					Windows: cfgWindows,
				})
				if err := m.config.Save(); err != nil {
					m.logger.Printf("Error saving config: %v", err)
					statusCmd := m.sessionList.NewStatusMessage(errorStyle("󰗼 save failed"))
					return m, statusCmd
				}
				statusCmd := m.sessionList.NewStatusMessage(statusMessageStyle("󰗼 saved to config"))
				return m, tea.Sequence(statusCmd, getSessionsCmd(m.config))
			case key.Matches(msg, m.keys.Start):
				si := m.sessionList.SelectedItem().(session)
				for _, s := range m.config.Sessions {
					if s.Name == si.sessionName {
						err := tmux.StartSession(s.Name, s.Directory, s.Windows)
						if err != nil {
							m.logger.Printf("Error starting session %s: %v", s.Name, err)
						}
						return m, tea.Quit
					}
				}
				return m, nil
			case key.Matches(msg, m.keys.Enter):
				si := m.sessionList.SelectedItem().(session)
				if si.sessionName == m.activeSession || !si.isRunning {
					statusCmd := m.sessionList.NewStatusMessage(errorStyle("󰗼 active or not running"))
					return m, statusCmd
				}
				if err := tmux.SwitchSession(si.sessionName); err != nil {
					m.logger.Printf("Error switching to session %s: %v", si.sessionName, err)
				}
				return m, tea.Quit
			case key.Matches(msg, m.keys.Help):
				m.help.ShowAll = !m.help.ShowAll
			case key.Matches(msg, defaultKeyMap.Quit):
				return m, tea.Quit
			}
		}
	case sessionsResponseMsg:
		m.activeSession = msg.ActiveSession
		m.sessionList.SetItems(msg.SessionsList)
	}
	m.sessionList, cmd = m.sessionList.Update(msg)
	cmds = append(cmds, cmd)
	m.sessionInput, cmd = m.sessionInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
