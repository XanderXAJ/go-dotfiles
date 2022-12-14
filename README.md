# Go Dotfiles

A utility to help with loading dotfiles in the terminal.
Currently `bash` only.

## Usage

1. Create a `.dotfiles/.bashrc.d` directory for your shell configuration.
2. Create as few or many shell configuration files inside as you wish, naming them in an `init.d` style (e.g. `00_system`, `10_git` etc.) to order them.
3. Add the follow to your start-up `rc` file (e.g. `bashrc`) to import your ordered correctly dotfiles:

```bash
for dotfile in $(go-dotfiles); do source "$dotfile"; done
```
