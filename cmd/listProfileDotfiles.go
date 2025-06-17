package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
)

type listProfileDotfilesConfig struct {
	shell string
}

var cliConfig listProfileDotfilesConfig

// listProfileDotfilesCmd represents the listProfileDotfiles command
var listProfileDotfilesCmd = &cobra.Command{
	Use:   "listProfileDotfiles",
	Short: "List ordered profile files in your dotfiles directories",
	Long: `Lists ordered profile files in your dotfile directories, ready for sourcing.
e.g. 00_system, 10_git, 99_zoxide etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		listDotfiles()
	},
}

func init() {
	rootCmd.AddCommand(listProfileDotfilesCmd)
	listProfileDotfilesCmd.Flags().StringVar(&cliConfig.shell, "shell", "", "Specify shell (e.g. bash, zsh, fish)")
}

func listDotfiles() {
	// TODO Get environment configuration

	var shell string
	if cliConfig.shell != "" {
		shell = cliConfig.shell
	} else {
		shell = detectShell()
	}

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

func detectShell() string {
	shell := "bash"
	if shellEnv := os.Getenv("SHELL"); shellEnv != "" {
		shell = filepath.Base(shellEnv)
	}
	if len(shell) > 0 && shell[0] == '.' {
		shell = shell[1:]
	}
	return shell
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
