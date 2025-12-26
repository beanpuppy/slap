//soppo:generated v1
// Help text generation for slap CLI parser.
package slap

import "reflect"
import "strings"

func formatHelp(meta *commandMeta) string {
	var b strings.Builder

	// Usage line
	b.WriteString("Usage: ")
	b.WriteString(meta.name)

	if len(meta.flags) > 0 {
		b.WriteString(" [OPTIONS]")
	}

	for _, am := range meta.args {
		if am.last {
			if am.required {
				b.WriteString(" <")
				b.WriteString(am.valueName)
				b.WriteString(">...")
			} else {
				b.WriteString(" [")
				b.WriteString(am.valueName)
				b.WriteString("]...")
			}
		} else {
			if am.required {
				b.WriteString(" <")
				b.WriteString(am.valueName)
				b.WriteString(">")
			} else {
				b.WriteString(" [")
				b.WriteString(am.valueName)
				b.WriteString("]")
			}
		}
	}

	b.WriteString("\n")

	// Description
	if meta.about != "" {
		b.WriteString("\n")
		b.WriteString(meta.about)
		b.WriteString("\n")
	}

	// Arguments section
	if len(meta.args) > 0 {
		b.WriteString("\nArguments:\n")
		for _, am := range meta.args {
			b.WriteString("  <")
			b.WriteString(am.valueName)
			b.WriteString(">")
			if am.help != "" {
				b.WriteString("  ")
				b.WriteString(am.help)
			}
			b.WriteString("\n")
		}
	}

	// Options section
	if len(meta.flags) > 0 {
		b.WriteString("\nOptions:\n")
		for _, fm := range meta.flags {
			if fm.hidden {
				continue
			}

			b.WriteString("  ")

			if fm.short != "" {
				b.WriteString("-")
				b.WriteString(fm.short)
				if fm.long != "" {
					b.WriteString(", ")
				}
			} else {
				b.WriteString("    ")
			}

			if fm.long != "" {
				b.WriteString("--")
				b.WriteString(fm.long)
			}

			// Add value placeholder for non-bool
			var fmKind reflect.Kind = fm.fieldType.Kind()
			if fmKind != reflect.Bool {
				b.WriteString(" <")
				b.WriteString(strings.ToUpper(fm.fieldName))
				b.WriteString(">")
			}

			if fm.help != "" {
				b.WriteString("  ")
				b.WriteString(fm.help)
			}

			// Show default/env
			if fm.defValue != "" || fm.env != "" {
				b.WriteString(" [")
				if fm.defValue != "" {
					b.WriteString("default: ")
					b.WriteString(fm.defValue)
				}
				if fm.defValue != "" && fm.env != "" {
					b.WriteString(", ")
				}
				if fm.env != "" {
					b.WriteString("env: ")
					b.WriteString(fm.env)
				}
				b.WriteString("]")
			}

			b.WriteString("\n")
		}

		// Always add help flag
		b.WriteString("  -h, --help  Print help\n")

		// Add version if available
		if meta.version != "" {
			b.WriteString("  -V, --version  Print version\n")
		}
	}

	return b.String()
}

// FormatHelp returns the help text for command type T.
func FormatHelp[T any]() string {
	var cmd T
	cmdType := reflect.TypeOf(cmd)
	target := getTarget(cmdType)
	meta := buildMeta(target, cmdType)
	return formatHelp((&meta))
}

