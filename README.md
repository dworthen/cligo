# cligo

A small, dependency-free Go library for building command-line applications. Flags are declared as a plain Go struct with `cligo` struct tags and are parsed into that struct with full type safety via generics.

## Features

- **Typed flags** – declare flags as struct fields; your `Run` function receives the populated struct.
- **Short flags, descriptions, required flags, and defaults** – all configured via struct tags.
- **Environment variable fallback** – any flag can be sourced from an environment variable.
- **`.env` file loading** – optionally load one or more dotenv files before parsing.
- **Positional arguments** – declare required (`<name>`) and optional (`[name]`) arguments.
- **Subcommands and aliases** – nest commands to any depth.
- **Generated help and version output** – `--help, -h` and `--version, -v` out of the box.
- **Strict mode** – reject unknown flags and unexpected arguments (enabled by default).

## Installation

```sh
go get github.com/dworthen/cligo
```

## Quick start

```go
package main

import (
	"fmt"

	"github.com/dworthen/cligo"
)

type greetFlags struct {
	Name  string `cligo:"required,short=n,env=GREET_NAME,description=who to greet"`
	Times int    `cligo:"short=t,description=number of greetings"`
}

func main() {
	cmd, err := cligo.NewCommand(
		cligo.WithName[greetFlags]("greet"),
		cligo.WithDescription[greetFlags]("Print a greeting"),
		cligo.WithFlags(greetFlags{Times: 1}), // default values
		cligo.WithRun(func(args []string, flags greetFlags) error {
			for i := 0; i < flags.Times; i++ {
				fmt.Printf("Hello, %s!\n", flags.Name)
			}
			return nil
		}),
	)
	if err != nil {
		panic(err)
	}

	cligo.MustRunCommand(cmd)
}
```

```sh
$ greet --name World -t 2
Hello, World!
Hello, World!

$ GREET_NAME=Gopher greet
Hello, Gopher!
```

## Defining commands

Commands are created with `NewCommand[T]`, where `T` is the flags struct type, using functional options:

| Option                          | Description                                                          |
| ------------------------------- | -------------------------------------------------------------------- |
| `WithName[T](name)`             | **Required.** Command name, shown in help and used to match subcommands. |
| `WithDescription[T](desc)`      | Description shown in help output.                                    |
| `WithAliases[T]([]string{...})` | Alternative names that can be used to invoke a subcommand.           |
| `WithArgs[T]("<a> [b]")`        | Positional argument specification (see [Positional arguments](#positional-arguments)). |
| `WithFlags(flags)`              | Flags struct; non-zero field values act as defaults. Must be a struct. |
| `WithRun(fn)`                   | `func(args []string, flags T) error` invoked when the command runs.  |
| `WithSubCommands[T](cmds...)`   | Nested subcommands.                                                  |

Commands that take no flags can use `struct{}` as `T`. Go cannot infer `T` for options such as `WithName` and `WithDescription`, so specify it explicitly. `WithFlags` and `WithRun` infer it from their arguments.

## Flags

Only struct fields with a `cligo` tag are treated as flags. The tag is a comma-separated list of the following entries:

| Tag entry            | Description                                                                 |
| -------------------- | --------------------------------------------------------------------------- |
| `required`           | The flag must be provided (via CLI, environment, or default). Ignored for `bool`. |
| `flag=<name>`        | Override the long flag name. Defaults to the hyphenated field name (`MaxRetries` → `--max-retries`). |
| `short=<c>`          | Short flag name, e.g. `short=n` → `-n`.                                     |
| `description=<text>` | Description shown in help output. Must not contain commas.                 |
| `env=<VAR>`          | Read the value from environment variable `VAR` when the flag isn't passed. |

Example:

```go
type flags struct {
	Output  string   `cligo:"required,short=o,description=output file"`
	Tags    []string `cligo:"flag=tag,short=t,env=APP_TAGS,description=tags to apply"`
	Verbose bool     `cligo:"short=v,description=enable verbose logging"`
	Ignored string   // no cligo tag; not a flag
}
```

### Supported types

`string`, `int`, `uint`, `float32`, `float64`, `bool`, and slices of each (`[]string`, `[]int`, …).

- **Value flags** take the next argument as their value: `--output out.txt` or `-o out.txt`.
- **Boolean flags** are set to `true` when present: `--verbose` or `-v`.
- **Slice flags** collect every occurrence: `-t a -t b` → `[]string{"a", "b"}`.

### Value precedence

1. Values passed on the command line.
2. The environment variable named by `env=`. For slice flags, the value is split on commas (`APP_TAGS=a,b`). For `bool` flags, `true` (case-insensitive) means true.
3. The default value set on the struct passed to `WithFlags`.

A `required` flag is satisfied by any of these sources. In practice, it fails only when the final value is still the zero value.

## Positional arguments

Declare positional arguments with `WithArgs`, using a space-separated spec:

- `<name>` – required argument
- `[name]` – optional argument

Required arguments cannot follow optional ones. Positional arguments must come **before** any flags. Parsing stops at the first argument that begins with `-`.

```go
cligo.WithArgs[flags]("<source> [destination]")
```

```sh
$ app src.txt dst.txt --verbose
```

The arguments are passed to your `Run` function as `args []string`. In strict mode, supplying more arguments than declared returns an error.

## Subcommands

```go
listCmd, _ := cligo.NewCommand(
	cligo.WithName[struct{}]("list"),
	cligo.WithAliases[struct{}]([]string{"ls"}),
	cligo.WithDescription[struct{}]("List items"),
	cligo.WithRun(func(args []string, _ struct{}) error {
		fmt.Println("listing...")
		return nil
	}),
)

rootCmd, _ := cligo.NewCommand(
	cligo.WithName[struct{}]("app"),
	cligo.WithDescription[struct{}]("My application"),
	cligo.WithSubCommands[struct{}](listCmd),
)

cligo.MustRunCommand(rootCmd)
```

```sh
$ app list
$ app ls        # alias
```

Each subcommand has its own flags type, arguments, and help output. The root command's name isn't part of the arguments; it only appears in help output.

## Running commands

| Function                       | Description                                                                 |
| ------------------------------ | --------------------------------------------------------------------------- |
| `RunCommand(cmd, opts...)`     | Run the command. Returns the command that was resolved (root or subcommand) and any error. |
| `MustRunCommand(cmd, opts...)` | Run the command. If an error occurs, print it with the command's usage to stderr and exit with status 1. |

### Options

Both functions accept optional functional options:

```go
cligo.MustRunCommand(cmd,
	cligo.WithVersion("1.2.3"),
	cligo.WithEnvFiles(".env", ".env.local"),
	cligo.WithStrict(true),
)
```

| Option                     | Default           | Description                                                      |
| -------------------------- | ----------------- | ---------------------------------------------------------------- |
| `WithVersion(v)`           | `"0"`             | Version string printed by the version flag.                      |
| `WithVersionFlag(flags)`   | `"--version, -v"` | Comma-separated flags that print the version. `""` disables it.  |
| `WithHelpFlag(flags)`      | `"--help, -h"`    | Comma-separated flags that print help. `""` disables it.         |
| `WithStrict(bool)`         | `true`            | Error on unknown flags and on extra positional arguments.        |
| `WithEnvFiles(paths...)`   | none              | Dotenv files to load before parsing.                             |
| `WithCliArgs(args...)`     | `os.Args[1:]`     | Command-line arguments to parse.                                 |

Use `WithCliArgs` to run a command against a custom argument list, such as in tests:

```go
_, err := cligo.RunCommand(rootCmd, cligo.WithCliArgs("list", "--verbose"))
```

### Env files

Files passed to `WithEnvFiles` are loaded once per process, in order, before any flags are parsed. Relative paths resolve against the current working directory, and missing files are skipped silently. The supported format is simple:

```sh
# comments and blank lines are ignored
APP_TAGS=a,b,c
GREETING=Hello $USER   # $VAR / ${VAR} references are expanded
```

Values loaded from env files overwrite variables already set in the environment.

## Help output

Usage text is generated from the command definition. For example, the sample CLI in [`internal/testingcli`](internal/testingcli/main.go) produces:

```text
start
	Start the entity

Aliases: run, init

Usage:
	start <name> [age] [...options]
	start <command>

Commands:
  cool    Cool command

Options:
  --name, -n (required) [default=$TESTING_NAME]    name of the entity
  --age, -a (required) [default=30]                age of the entity
  --files, -f [default=[file1.txt file2.txt]]      list of files
  --olde-enough, -o                                whether the entity is old enough
  --score, -s [default=$TESTING_SCORE]             score of the entity
  --help, -h                                       Show help
  --version, -v                                    Show version
```

## Development

This project uses [Task](https://taskfile.dev), which is declared as a Go tool dependency in `go.mod`, so no separate install is needed:

```sh
go tool task tidy          # go mod tidy
go tool task format        # goimports -w .
go tool task format-check  # list files needing formatting
go tool task lint          # golangci-lint run
go tool task run -- --help # run the sample CLI in internal/testingcli
```

`format`, `format-check`, and `lint` require [`goimports`](https://pkg.go.dev/golang.org/x/tools/cmd/goimports) and [`golangci-lint`](https://golangci-lint.run) on your `PATH`. The Taskfile loads variables from a local `.env` file, which is git-ignored.
