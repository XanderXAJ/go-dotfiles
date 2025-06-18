package discovery

import (
	"os"
	"path/filepath"
)

func DotfileDirs(dir string) ([]string, error) {
	var basePaths []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			matched, _ := filepath.Match(".dotfiles*", entry.Name())
			if matched {
				basePaths = append(basePaths, filepath.Join(dir, entry.Name()))
			}
		}
	}
	return basePaths, nil
}
