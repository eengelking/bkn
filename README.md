# bkn
![Bacon](images/bacon.png)

Pronounced as Bacon. An attempt to create a simplified Make and Makefile for running commands.

This is a simple program that allows you to run commands from a YAML file. This is useful for running commands that you run often but don't want to remember the syntax for, especially in pipelines.

## Binary Requirements
- AMD64 Linux

## Binary Installation
1. CD to the `bin` directory.
2. Run `./bkn` from this directory or copy the `bkn` binary and the `bkn.yaml` to another directory.
3. Update the bkn.yaml file with your commands.

## Building from Source
1. Install Golang
2. Clone the repo
3. CD to the repo directory
4. Run `go build -o bin/bkn main.go`

## YAML File
The YAML file is used to define the commands that you want to run. The file should be named `bkn.yaml` and should be in the same directory as the `bkn` binary.

Example:
```yaml
include:  # Optional. Include other YAML files to keep your commands organized.
  - /path/to/other/file.yaml

commands:
  - name: list  # The name of the command
    description: List the contents of the directory.  # The description of the command
    command: |- # The command to run. The |- allows for multiline commands in a YAML file
      ls -lah "$@"

  - name: listening
    description: List all programs listening on a given port
    command: |-
      ss -tupln | grep LISTEN
```

### Passing Variables
Any arguments after the command name are forwarded to the shell as positional parameters, just like a bash script. Reference them in your YAML command with `$1`, `$2`, `$3`, …, or use `$@` to expand all of them.

For example, given the `list` command above:

- `./bkn list` runs `ls -lah` with no arguments.
- `./bkn list /tmp` runs `ls -lah /tmp`.
- `./bkn list /tmp /var` runs `ls -lah /tmp /var`.

You can also pull individual positions, e.g. `command: echo "first=$1 second=$2"`.

# Attributes
* <a href="https://www.flaticon.com/free-icons/bacon" title="bacon icons">Bacon icons created by Freepik - Flaticon</a>