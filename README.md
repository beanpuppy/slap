# slap

A clap-style CLI parser for [Soppo](https://github.com/halcyonnouveau/soppo).

## Example

```go
package main

import (
	"os"
	
	slap "github.com/beanpuppy/slap/gen"
)

[slap.Command{Name: "add", About: "Add a task"}]
type AddCmd struct {
	[slap.Arg{Position: 0, Help: "Task description"}]
	Task string
}

func (cmd AddCmd) Run() error {
	fmt.Println("Added: {cmd.Task}")
	return nil
}

[slap.Command{Name: "todo", About: "Manage tasks"}]
[slap.Subcommands]
type Cmd enum {
	Add AddCmd

	[slap.Alias{Name: "ls"}]
	List ListCmd

	[slap.Hidden]
	Debug DebugCmd
}

func main() {
	slap.Run[Cmd]() ? {
		os.Exit(1)
	}
}
```

```
$ todo
Usage: todo [OPTIONS] <COMMAND>

Manage tasks

Commands:
  add       Add a task
  list (ls) List tasks

Options:
  -h, --help  Print help
```

## Attributes

| Attribute | Fields |
|-----------|--------|
| `Command` | `Name`, `About`, `Version` |
| `Subcommands` | Marks enum as subcommand container |
| `Arg` | `Position`, `Help`, `Required`, `ValueName`, `Last` |
| `Flag` | `Short`, `Long`, `Help`, `Default`, `Env`, `Hidden`, `Global` |
| `Hidden` | Hide subcommand from help |
| `Alias` | `Name` - alternative subcommand name |

## License

BSD 3-Clause
