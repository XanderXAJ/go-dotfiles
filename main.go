package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	// TODO Get environment configuration

	// TODO Get shell
	shell := "bash"

	// TODO Get dotfiles basePath/matcher
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	basePaths, err := dotfileDirs(home)
	if err != nil {
		log.Fatalln("Error getting dotfile directories:", err)
	}

	// Discover dotfiles
	dotfiles := make([]string, 0, 50)
	for _, basePath := range basePaths {
		filepath.WalkDir(basePath, func(walkPath string, d fs.DirEntry, err error) error {
			log.Println(walkPath, err)
			if walkPath == basePath {
				return nil
			}
			if d.IsDir() {
				return fs.SkipDir
			}
			if d.Type().IsRegular() && filepath.Ext(walkPath) == "" || filepath.Ext((walkPath)) == fmt.Sprintf(".%s", shell) {
				dotfiles = append(dotfiles, walkPath)
			}
			return nil
		})
	}

	// Sort dotfiles in to order
	sort.Slice(dotfiles, func(i, j int) bool {
		return filepath.Base(dotfiles[i]) < filepath.Base(dotfiles[j])
	})

	// Print dotfiles
	for _, dotfile := range dotfiles {
		fmt.Println(dotfile)
	}
}

func dotfileDirs(dir string) ([]string, error) {
	var basePaths []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			matched, _ := filepath.Match(".dotfiles*", entry.Name())
			if matched {
				basePaths = append(basePaths, filepath.Join(dir, entry.Name(), ".bashrc.d"))
			}
		}
	}
	return basePaths, nil
}
