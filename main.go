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
	basePaths := []string{
		filepath.Join(home, ".dotfiles/.bashrc.d"),
		filepath.Join(home, ".dotfiles_test/.bashrc.d"),
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
