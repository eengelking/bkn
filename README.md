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
      ls -lah <variable>

  - name: listening
    description: List all programs listening on a given port
    command: |-
      ss -tupln | grep LISTEN
```

### The Purpose of \<VARIALBE>
You can pass variables with the `./bkn` command.

For example, if I run the list command from the example above, I would run `./bkn list`. If I wanted to pass a variable to the command, I would run `./bkn list <variable>`. The variable is then replaced in the command with the value that you passed. For example, if I ran `./bkn list /tmp`, the command would be `ls -lah /tmp`.

At the moment, you can only pass one variable to a command.

# Attributes
* <a href="https://www.flaticon.com/free-icons/bacon" title="bacon icons">Bacon icons created by Freepik - Flaticon</a>