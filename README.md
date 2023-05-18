# bkn
Pronounced as Bacon. A better Make using Golang

## Installation
1. Copy the `bkn` binary to your directory.
2. Update the bkn.yaml file with your commands.

## YAML File
The YAML file is used to define the commands that you want to run. The file should be named `bkn.yaml` and should be in the same directory as the `bkn` binary.

Example:
```yaml
- name: list  # The name of the command
  description: List the contents of the directory.  # The description of the command
  command: |- # The command to run. The |- allows for multiline commands in a YAML file
    ls -lah <variable>

- name: listening
  description: List all programs listening on a given port
  command: |-
    ss -tupln | grep LISTEN
```

## The Purpose of <VARIALBE>
You can pass variables with the `./bkn` command.

For example, if I run the list command from the example above, I would run `./bkn list`. If I wanted to pass a variable to the command, I would run `./bkn list <variable>`. The variable is then replaced in the command with the value that you passed. For example, if I ran `./bkn list /tmp`, the command would be `ls -lah /tmp`.

At the moment, you can only pass one variable to a command.