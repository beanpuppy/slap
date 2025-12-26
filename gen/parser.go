//soppo:generated v1
// Core parsing logic for slap CLI parser.
package slap

import "fmt"
import "os"
import "reflect"
import "strconv"
import "strings"

// Parse parses os.Args[1:] into a command enum S.
// If S has only one variant, no subcommand name is required.
func Parse[S any]() (S, error) {
	return ParseArgs[S](os.Args[1:])
}

// emptyParent is used when parsing without a parent struct
type emptyParent struct {
}

// ParseArgs parses the given arguments into a command enum S.
func ParseArgs[S any](args []string) (S, error) {
	var sub S
	_, sub, err := ParseArgsSub[emptyParent, S](args)
	return sub, err
}

// Run parses os.Args and executes the matched command.
// Each variant in S must implement Runnable.
func Run[S any]() error {
	sub, _err0 := Parse[S]()
	if _err0 != nil {
		err := _err0
		var parseErr ParseError
		if ok := errorAs(err, (&parseErr)); ok {
			switch parseErr.Kind.(type) {
			case ErrorKind_HelpRequested:
				fmt.Println(parseErr.Message)
				return nil
			case ErrorKind_VersionRequested:
				fmt.Println(parseErr.Message)
				return nil
			default:
			}
		}
		// Fall through
		return err
	}
	return runEnumVariant(sub)
}

// runEnumVariant extracts and runs the active variant from an enum
func runEnumVariant(sub any) error {
	// Try direct Runnable (unlikely for enum interface)
	if runner, ok := sub.(Runnable); ok {
		return runner.Run()
	}

	subVal := reflect.ValueOf(sub)

	// For interface-based enums, the value is the variant wrapper struct
	// Get the underlying concrete value if it's an interface
	if subVal.Kind() == reflect.Interface {
		subVal = subVal.Elem()
	}

	// The wrapper struct has a Value field containing the command
	if subVal.Kind() == reflect.Struct {
		valueField := subVal.FieldByName("Value")
		if valueField.IsValid() {
			if runner, ok := valueField.Interface().(Runnable); ok {
				return runner.Run()
			}
		}
	}

	return fmt.Errorf("command does not implement Runnable")
}

func parseInto(cmdVal reflect.Value, meta *commandMeta, args []string) error {
	// Track which flags and args have been set
	flagsSet := make(map[string]bool)
	argsSet := make(map[int]bool)

	// Collect positional args (non-flag arguments)
	positionals := []string{}
	afterDoubleDash := false

	i := 0
	for i < len(args) {
		arg := args[i]

		// Handle -- separator
		if arg == "--" {
			afterDoubleDash = true
			i++
			continue
		}

		// After --, everything is positional
		if afterDoubleDash {
			positionals = append(positionals, arg)
			i++
			continue
		}

		// Handle flags
		if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			var flagName string
			var value string
			hasValue := false

			if strings.HasPrefix(arg, "--") {
				// Long flag
				flagPart := arg[2:]
				if eqIdx := strings.Index(flagPart, "="); eqIdx >= 0 {
					flagName = flagPart[:eqIdx]
					value = flagPart[eqIdx + 1:]
					hasValue = true
				} else {
					flagName = flagPart
				}
			} else {
				// Short flag
				flagName = arg[1:2]
				if len(arg) > 2 {
					// -fvalue or -f=value
					rest := arg[2:]
					if rest[0] == '=' {
						value = rest[1:]
					} else {
						value = rest
					}
					hasValue = true
				}
			}

			// Find matching flag
			var fm *flagMeta = nil
			for j := range meta.flags {
				if meta.flags[j].short == flagName || meta.flags[j].long == flagName {
					fm = (&meta.flags[j])
					break
				}
			}

			if fm == nil {
				return ParseError{
					Kind: ErrorKindUnknownFlag,
					Value: arg,
					Message: "unknown flag: " + arg,
				}
			}

			// Get value if needed
			var fieldKind reflect.Kind = fm.fieldType.Kind()
			if fieldKind == reflect.Bool {
				// Boolean flags don't need a value
				if (!hasValue) {
					value = "true"
				}
			} else {
				if (!hasValue) {
					// Need next arg as value
					i++
					if i >= len(args) {
						return ParseError{
							Kind: ErrorKindInvalidValue,
							Field: fm.fieldName,
							Message: "flag --" + fm.long + " requires a value",
						}
					}
					value = args[i]
				}
			}

			// Set the field
			field := cmdVal.Field(fm.fieldIdx)
			_err0 := setField(field, value, fm.fieldType)
			if _err0 != nil {
				err := _err0
				return ParseError{
					Kind: ErrorKindInvalidValue,
					Field: fm.fieldName,
					Value: value,
					Message: "invalid value '" + value + "' for --" + fm.long + ": " + err.Error(),
				}
			}
			flagsSet[fm.fieldName] = true
		} else {
			// Positional argument
			positionals = append(positionals, arg)
		}

		i++
	}

	// Assign positional arguments
	posIdx := 0
	for _, am := range meta.args {
		if am.last {
			// Consume all remaining positionals
			if posIdx < len(positionals) {
				field := cmdVal.Field(am.fieldIdx)
				remaining := positionals[posIdx:]
				_err1 := setSliceField(field, remaining, am.fieldType)
				if _err1 != nil {
					err := _err1
					return ParseError{
						Kind: ErrorKindInvalidValue,
						Field: am.fieldName,
						Message: err.Error(),
					}
				}
				posIdx = len(positionals)
			}
			argsSet[am.position] = true
		} else {
			if posIdx < len(positionals) {
				field := cmdVal.Field(am.fieldIdx)
				_err2 := setField(field, positionals[posIdx], am.fieldType)
				if _err2 != nil {
					err := _err2
					return ParseError{
						Kind: ErrorKindInvalidValue,
						Field: am.fieldName,
						Value: positionals[posIdx],
						Message: "invalid value for <" + am.valueName + ">: " + err.Error(),
					}
				}
				argsSet[am.position] = true
				posIdx++
			}
		}
	}

	// Check for required args
	for _, am := range meta.args {
		if am.required && (!argsSet[am.position]) {
			return ParseError{
				Kind: ErrorKindMissingRequired,
				Field: am.fieldName,
				Message: "missing required argument: <" + am.valueName + ">",
			}
		}
	}

	// Apply defaults and env vars for unset flags
	for _, fm := range meta.flags {
		if flagsSet[fm.fieldName] {
			continue
		}

		// Try environment variable
		if fm.env != "" {
			if envVal, ok := os.LookupEnv(fm.env); ok {
				field := cmdVal.Field(fm.fieldIdx)
				if err := setField(field, envVal, fm.fieldType); err == nil {
					continue
				}
			}
		}

		// Try default
		if fm.defValue != "" {
			field := cmdVal.Field(fm.fieldIdx)
			_ = setField(field, fm.defValue, fm.fieldType)
		}
	}

	return nil
}

func setField(field reflect.Value, value string, t reflect.Type) error {
	var kind reflect.Kind = t.Kind()
	switch kind {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, _err0 := strconv.ParseInt(value, 10, 64)
		if _err0 != nil {
			return _err0
		}
		field.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, _err1 := strconv.ParseUint(value, 10, 64)
		if _err1 != nil {
			return _err1
		}
		field.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, _err2 := strconv.ParseFloat(value, 64)
		if _err2 != nil {
			return _err2
		}
		field.SetFloat(n)
	case reflect.Bool:
		b, _err3 := strconv.ParseBool(value)
		if _err3 != nil {
			return _err3
		}
		field.SetBool(b)
	default:
		return fmt.Errorf("unsupported type: %s", kind)
	}
	return nil
}

func setSliceField(field reflect.Value, values []string, t reflect.Type) error {
	var kind reflect.Kind = t.Kind()
	if kind != reflect.Slice {
		return fmt.Errorf("expected slice type, got %s", kind)
	}

	elemType := t.Elem()
	slice := reflect.MakeSlice(t, len(values), len(values))

	for i, v := range values {
		elem := slice.Index(i)
		_err0 := setField(elem, v, elemType)
		if _err0 != nil {
			return _err0
		}
	}

	field.Set(slice)
	return nil
}

