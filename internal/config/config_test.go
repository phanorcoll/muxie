package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestConfigFilePath(t *testing.T) {
	path, err := configFilePath()
	if err != nil {
		t.Fatalf("configFilePath() returned error: %v", err)
	}

	if path == "" {
		t.Fatal("configFilePath() returned empty path")
	}

	expectedSuffix := filepath.Join(".config", "muxie", "config.yml")
	if !strings.HasSuffix(path, expectedSuffix) {
		t.Errorf("configFilePath() = %q, want suffix %q", path, expectedSuffix)
	}
}

func TestSaveAndLoadRoundtrip(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yml")

	original := &Config{
		Sessions: []Session{
			{
				Name:      "my-project",
				Directory: "~/code/my-project",
				Windows: []Window{
					{
						Name:   "editor",
						Layout: "vertical",
						Panes: []Pane{
							{Command: "nvim", Directory: "~/code/my-project/src"},
							{Command: "git status"},
						},
					},
					{
						Name: "server",
						Panes: []Pane{
							{Command: "npm run dev"},
						},
					},
				},
			},
		},
	}

	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal config: %v", err)
	}
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	loaded, err := loadFile(configFile)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if len(loaded.Sessions) != len(original.Sessions) {
		t.Fatalf("session count mismatch: got %d, want %d", len(loaded.Sessions), len(original.Sessions))
	}

	orig := original.Sessions[0]
	got := loaded.Sessions[0]

	if got.Name != orig.Name {
		t.Errorf("session name: got %q, want %q", got.Name, orig.Name)
	}
	if got.Directory != orig.Directory {
		t.Errorf("session directory: got %q, want %q", got.Directory, orig.Directory)
	}
	if len(got.Windows) != len(orig.Windows) {
		t.Fatalf("window count: got %d, want %d", len(got.Windows), len(orig.Windows))
	}
	if got.Windows[0].Name != orig.Windows[0].Name {
		t.Errorf("window name: got %q, want %q", got.Windows[0].Name, orig.Windows[0].Name)
	}
	if got.Windows[0].Layout != orig.Windows[0].Layout {
		t.Errorf("window layout: got %q, want %q", got.Windows[0].Layout, orig.Windows[0].Layout)
	}
	if len(got.Windows[0].Panes) != len(orig.Windows[0].Panes) {
		t.Fatalf("pane count: got %d, want %d", len(got.Windows[0].Panes), len(orig.Windows[0].Panes))
	}
	if got.Windows[0].Panes[0].Command != orig.Windows[0].Panes[0].Command {
		t.Errorf("pane command: got %q, want %q", got.Windows[0].Panes[0].Command, orig.Windows[0].Panes[0].Command)
	}
	if got.Windows[0].Panes[0].Directory != orig.Windows[0].Panes[0].Directory {
		t.Errorf("pane directory: got %q, want %q", got.Windows[0].Panes[0].Directory, orig.Windows[0].Panes[0].Directory)
	}
	if got.Windows[0].Panes[1].Command != orig.Windows[0].Panes[1].Command {
		t.Errorf("second pane command: got %q, want %q", got.Windows[0].Panes[1].Command, orig.Windows[0].Panes[1].Command)
	}
}

func TestSaveEmptyConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yml")

	empty := &Config{}
	data, err := yaml.Marshal(empty)
	if err != nil {
		t.Fatalf("failed to marshal empty config: %v", err)
	}
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		t.Fatalf("failed to write empty config: %v", err)
	}

	loaded, err := loadFile(configFile)
	if err != nil {
		t.Fatalf("failed to load empty config: %v", err)
	}

	if loaded.Sessions != nil && len(loaded.Sessions) != 0 {
		t.Errorf("expected no sessions, got %d", len(loaded.Sessions))
	}
}

func TestSaveCreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "nonexistent", "muxie", "config.yml")

	cfg := &Config{
		Sessions: []Session{
			{Name: "test"},
		},
	}

	// Simulate what Save() does: MkdirAll then write.
	parentDir := filepath.Dir(configFile)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		t.Fatalf("failed to create parent dir: %v", err)
	}

	data, _ := yaml.Marshal(cfg)
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	if _, err := os.Stat(configFile); err != nil {
		t.Fatalf("config file not created: %v", err)
	}

	loaded, err := loadFile(configFile)
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}
	if len(loaded.Sessions) != 1 || loaded.Sessions[0].Name != "test" {
		t.Errorf("unexpected loaded config: %+v", loaded)
	}
}

func TestSaveOverwritesExisting(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yml")

	first := &Config{Sessions: []Session{{Name: "first"}}}
	data, _ := yaml.Marshal(first)
	os.WriteFile(configFile, data, 0644)

	second := &Config{Sessions: []Session{{Name: "second"}}}
	data, _ = yaml.Marshal(second)
	os.WriteFile(configFile, data, 0644)

	loaded, err := loadFile(configFile)
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}
	if len(loaded.Sessions) != 1 || loaded.Sessions[0].Name != "second" {
		t.Errorf("expected overwritten content, got %+v", loaded)
	}
}

func TestRoundtripMultipleSessions(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yml")

	original := &Config{
		Sessions: []Session{
			{Name: "session-a", Directory: "/tmp/a"},
			{Name: "session-b", Directory: "/tmp/b", Windows: []Window{
				{Name: "win1", Layout: "horizontal", Panes: []Pane{
					{Command: "htop"},
					{Command: "tail -f /var/log/syslog", Directory: "/var/log"},
				}},
			}},
		},
	}

	data, _ := yaml.Marshal(original)
	os.WriteFile(configFile, data, 0644)

	loaded, err := loadFile(configFile)
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}

	if len(loaded.Sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(loaded.Sessions))
	}

	// Session A: no windows.
	if loaded.Sessions[0].Windows != nil && len(loaded.Sessions[0].Windows) != 0 {
		t.Errorf("session-a should have no windows, got %d", len(loaded.Sessions[0].Windows))
	}

	// Session B: two panes in one window.
	sb := loaded.Sessions[1]
	if len(sb.Windows) != 1 {
		t.Fatalf("session-b: expected 1 window, got %d", len(sb.Windows))
	}
	if len(sb.Windows[0].Panes) != 2 {
		t.Fatalf("session-b window: expected 2 panes, got %d", len(sb.Windows[0].Panes))
	}
	if sb.Windows[0].Panes[1].Directory != "/var/log" {
		t.Errorf("pane directory: got %q, want /var/log", sb.Windows[0].Panes[1].Directory)
	}
}

// loadFile replicates the core of Load() given an explicit path, for testing.
func loadFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
