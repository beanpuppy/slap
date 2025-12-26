//soppo:generated v1
// Attribute structs for slap CLI parser.
// Command marks a struct as a CLI command.
package slap

type Command struct {
	Name string
	About string
	Version string
}

// Command name (defaults to lowercase struct name)
// Short description
// Version string
// Subcommands marks an enum as a collection of CLI subcommands.
// This enables runtime discovery of enum variants for CLI parsing.
type Subcommands struct {
}

// Hidden marks a subcommand variant as hidden from help output.
// The subcommand can still be used, it just won't appear in --help.
type Hidden struct {
}

// Alias adds an alternative name for a subcommand variant.
// Multiple Alias attributes can be applied to add multiple aliases.
type Alias struct {
	Name string
}

// Alternative name for the subcommand
// Arg marks a field as a positional argument.
type Arg struct {
	Position int
	Name string
	Help string
	Required bool
	ValueName string
	Last bool
}

// 0-based position index
// Display name (defaults to field name)
// Help text
// Required? (default: true for non-pointer types)
// Placeholder in help, e.g. "<FILE>"
// Capture all remaining args (for []T fields)
// Flag marks a field as a CLI flag.
type Flag struct {
	Short string
	Long string
	Help string
	Default string
	Env string
	Hidden bool
	Global bool
}

// Short flag, e.g. "v" for -v
// Long flag (defaults to lowercase field name)
// Help text
// Default value (as string, coerced to field type)
// Environment variable fallback
// Hide from help
// Available in subcommands
// Possible restricts values to a set of options.
type Possible struct {
	Values []string
}

// Range restricts numeric values to a range.
type Range struct {
	Min float64
	Max float64
}

