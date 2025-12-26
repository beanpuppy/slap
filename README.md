# slap

A clap-style CLI parser for [Soppo](https://github.com/halcyonnouveau/soppo).

```bash
go get github.com/beanpuppy/slap
```

## Example

```go
package main

import (
	"fmt"
	"os"
	slap "github.com/beanpuppy/slap/gen"
)

[slap.Command{Name: "greet", About: "Greet someone", Version: "1.0.0"}]
type GreetCmd struct {
	[slap.Arg{Position: 0, Help: "Name to greet"}]
	Name string

	[slap.Flag{Short: "c", Long: "count", Default: "1"}]
	Count int

	[slap.Flag{Short: "l", Long: "loud"}]
	Loud bool
}

func (cmd GreetCmd) Run() error {
	fmt.Println("Hello, {cmd.Name}!")
	return nil
}

[slap.Subcommands]
type Cmd enum {
	Greet GreetCmd
}

func main() {
	slap.Run[Cmd]() ? err {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

## Attributes

### Type/Field Attributes

| Attribute | Fields |
|-----------|--------|
| `Command` | `Name`, `About`, `Version` |
| `Subcommands` | (none) - marks an enum as CLI subcommands |
| `Arg` | `Position`, `Name`, `Help`, `Required`, `ValueName`, `Last` (for `[]T`) |
| `Flag` | `Short`, `Long`, `Help`, `Default`, `Env`, `Hidden`, `Global` |

### Variant Attributes

| Attribute | Fields | Description |
|-----------|--------|-------------|
| `Hidden` | (none) | Hide subcommand from `--help` (still usable) |
| `Alias` | `Name` | Add an alternative name for the subcommand |

## Subcommands

```go
[slap.Command{Name: "add"}]
type AddCmd struct { ... }

[slap.Command{Name: "remove"}]
type RemoveCmd struct { ... }

[slap.Subcommands]
type Cmd enum {
	Add AddCmd

	[slap.Alias{Name: "rm"}]
	Remove RemoveCmd

	[slap.Hidden]
	Debug DebugCmd
}

func main() {
	// Simple: just the subcommand enum
	slap.Run[Cmd]() ?

	// Or with global flags on a parent struct
	app, sub := slap.ParseSub[App, Cmd]() ?
	match sub {
	case Cmd.Add(cmd): ...
	case Cmd.Remove(cmd): ...
	}
}

[slap.Command{Name: "app"}]
type App struct {
	[slap.Flag{Short: "v", Global: true}]
	Verbose bool
}
```

## License

BSD 3-Clause. See [LICENSE](LICENSE).
