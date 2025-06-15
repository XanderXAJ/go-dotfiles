package cmd

import (
	"os"
	"testing"
)

func TestDetectShell(t *testing.T) {
	tests := []struct {
		name     string
		shellEnv string
		want     string
	}{
		{"empty SHELL", "", "bash"},
		{"absolute bash", "/bin/bash", "bash"},
		{"absolute zsh", "/usr/bin/zsh", "zsh"},
		{"dot prefix", ".bash", "bash"},
		{"dot prefix abs", "/foo/.bash", "bash"},
		{"fish", "/usr/local/bin/fish", "fish"},
		{"dash", "/usr/bin/dash", "dash"},
	}

	origShell := os.Getenv("SHELL")
	defer os.Setenv("SHELL", origShell)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shellEnv != "" {
				os.Setenv("SHELL", tt.shellEnv)
			} else {
				os.Unsetenv("SHELL")
			}
			got := detectShell()
			if got != tt.want {
				t.Errorf("detectShell() = %q, want %q", got, tt.want)
			}
		})
	}
}
