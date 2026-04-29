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

By default `bkn` reads `bkn.yaml` from the current working directory. Use `-c` / `--config` (see [Flags](#flags)) to point it at a YAML file in another location instead of copying one next to the binary.

## Flags

| Flag | Description |
| --- | --- |
| `-c`, `--config <path>` | Path to the YAML file to load. Defaults to `bkn.yaml` in the current working directory. |
| `-h`, `--help` | Print colored usage info and the list of available commands, then exit. |

Both short and long forms accept the standard Go-flag spellings: `-c path`, `-c=path`, `--config path`, and `--config=path` all work.

Examples:

```sh
# Run the `hello` command from a YAML file in another directory
bkn -c /etc/bkn/team.yaml hello

# Forward args to the underlying shell command (positional args after the command name)
bkn --config ~/dotfiles/bkn.yaml greet Ed

# Show help with the command list
bkn --help
```

Note: because flag parsing uses Go's standard `flag` package, `-c` / `-h` are consumed wherever they appear on the command line. If you need to forward a literal `-c` or `-h` to your shell command, quote it or rename the conflicting argument.

## Building from Source
1. Install Golang
2. Clone the repo
3. CD to the repo directory
4. Run `go build -o bin/bkn ./cmd/bkn`
5. Run `go test ./...` to execute the test suite

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

### Shell and syntax by OS
`bkn` invokes a different shell per host so positional-argument forwarding feels native on each platform:

| OS | Shell | Positional syntax |
| --- | --- | --- |
| Linux, macOS | `sh -c <body> bkn <args...>` | `$1`, `$2`, `$@`, `$#` |
| Windows | `cmd.exe /C <tempfile>.bat <args...>` (the body is written to a temp `.bat` so that `%1` substitution actually fires — `cmd /C` does not perform it on inline strings) | `%1`, `%2`, `%*` |

Command bodies are therefore **not portable** across OSes. Pair OS-specific bodies with the `os:` field so each host only sees the bodies it can run:

```yaml
commands:
  - name: list
    description: List a directory (linux/macOS).
    os: [linux, darwin]
    command: ls -lah "$@"

  - name: win-list
    description: List a directory (Windows).
    os: [windows]
    command: dir %1
```

### Restricting Commands by OS
Each command may declare an optional `os:` list to restrict it to specific host operating systems. Values match Go's `runtime.GOOS` (`linux`, `darwin`, `windows`, …). When `os:` is omitted, the command is available everywhere. When set, the command is only listed and only runnable on a matching host; on any other host it is hidden, and invoking it falls through to the standard `Invalid option` message.

```yaml
commands:
  - name: brew-update
    description: Update Homebrew (macOS only).
    os: [darwin]
    command: brew update && brew upgrade

  - name: ports
    description: List listening ports.
    os: [linux, darwin]
    command: ss -tupln | grep LISTEN
```

# Attributes
* <a href="https://www.flaticon.com/free-icons/bacon" title="bacon icons">Bacon icons created by Freepik - Flaticon</a>