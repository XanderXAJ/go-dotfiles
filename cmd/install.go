package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/xanderxaj/go-dotfiles/internal/discovery"
)

var dotfiles = []string{
	".bash_logout",
	".bash_profile",
	".bashrc",
	".gitconfig",
	".gitconfig_env",
	".gitconfig_private",
	".gitignore.global",
	".inputrc",
	".p10k.zsh",
	".profile",
	".tmux*.conf",
	".vimrc",
	".XCompose",
	".zshrc",
}

type installFlags struct {
	dryRun bool
}

var installRunConfig installFlags

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install dotfiles",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		install()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().BoolVarP(&installRunConfig.dryRun, "dry-run", "d", false, "Dry run (don't perform installation)")
}

func install() {
	// Get dotfiles directories
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}
	dotfileDirs, err := discovery.DotfileDirs(home)
	if err != nil {
		log.Fatalln(err)
	}

	for _, dotfile := range dotfiles {
		var sourcePath string

		// Find the first occurrence of the dotfile in the dotfileDirs
		for _, dir := range dotfileDirs {
			testPath := filepath.Join(dir, dotfile)
			fmt.Println("Testing", testPath)
			matches, err := filepath.Glob(testPath)
			if err != nil {
				log.Printf("Error globbing %s in %s: %v", dotfile, dir, err)
				continue
			}
			if len(matches) > 0 {
				// TODO: Support multiple glob matches
				sourcePath = matches[0]
				break
			}
		}
		if sourcePath == "" {
			continue // Not found in any dotfileDirs
		}
		fmt.Println("Found", sourcePath)

		targetPath := filepath.Join(home, filepath.Base(sourcePath))

		// Backup existing file
		if _, err := os.Lstat(targetPath); err == nil {
			backupPath := targetPath + ".bak." + time.Now().Format("20060102-150405")
			fmt.Println("Backing up", targetPath, "to", backupPath)
			if !installRunConfig.dryRun {
				err := os.Rename(targetPath, backupPath)
				if err != nil {
					log.Printf("Failed to backup %s: %v", targetPath, err)
					continue
				}
			}
		}
		// Symlink the file
		fmt.Println("Installing", targetPath)
		if !installRunConfig.dryRun {
			err := os.Symlink(sourcePath, targetPath)
			if err != nil {
				log.Printf("Failed to symlink %s to %s: %v", sourcePath, targetPath, err)
			}
		}
	}
}
