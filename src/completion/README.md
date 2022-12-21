# Exeiac completion
A small package to provide an auto-completion to different shells (bash, zsh, fish).

It justs references all brick's in a configuration file's "Rooms" section, and return them.

## Setup
*Instructions on how to build the binary used for brick's auto-completion.*
```bash
$ go install completion
```

## Bash

## Zsh

## Fish
To enable the fish auto-completion, you can just copy the `./scripts/exeiac.fish` file to `~/.config/fish/completions/exeiac.fish`.
```fish
$ cp ./scripts/exeiac.fish ~/.config/fish/completions/exeiac.fish
```
You should have completion enabled once
