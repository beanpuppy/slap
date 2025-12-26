//soppo:generated v1
// Tests for slap CLI parser.
package slap

import "github.com/halcyonnouveau/soppo/runtime"
import "strings"
import "testing"

// Command for basic parsing tests
type GreetTestCmd struct {
	Name string
	Count int
	Verbose bool
}

func (cmd GreetTestCmd) Run() error {
	return nil
}

/*soppo:enum
SingleCmd {
    Greet GreetTestCmd
}
*/
type SingleCmd interface {
	isSingleCmd()
}

type SingleCmd_Greet struct {
	Value GreetTestCmd
}
func (SingleCmd_Greet) isSingleCmd() {}

func SingleCmdGreet(value GreetTestCmd) SingleCmd {
	return SingleCmd_Greet{Value: value}
}

// Commands for subcommand tests
type AddItemCmd struct {
	Item string
}

func (cmd AddItemCmd) Run() error {
	return nil
}

type ListItemsCmd struct {
	All bool
}

func (cmd ListItemsCmd) Run() error {
	return nil
}

/*soppo:enum
MultiCmd {
    Add AddItemCmd
    List ListItemsCmd
}
*/
type MultiCmd interface {
	isMultiCmd()
}

type MultiCmd_Add struct {
	Value AddItemCmd
}
func (MultiCmd_Add) isMultiCmd() {}

type MultiCmd_List struct {
	Value ListItemsCmd
}
func (MultiCmd_List) isMultiCmd() {}

func MultiCmdAdd(value AddItemCmd) MultiCmd {
	return MultiCmd_Add{Value: value}
}
func MultiCmdList(value ListItemsCmd) MultiCmd {
	return MultiCmd_List{Value: value}
}

// Parent with global flags
type ParentCmd struct {
	Verbose bool
}

// Commands for Hidden/Alias tests
type CreateCmd struct {
	Name string
}

func (cmd CreateCmd) Run() error {
	return nil
}

type DeleteCmd struct {
	Name string
}

func (cmd DeleteCmd) Run() error {
	return nil
}

type DebugCmd struct {
}

func (cmd DebugCmd) Run() error {
	return nil
}

/*soppo:enum
AliasCmd {
    Create CreateCmd
    Delete DeleteCmd
    Debug DebugCmd
}
*/
type AliasCmd interface {
	isAliasCmd()
}

type AliasCmd_Create struct {
	Value CreateCmd
}
func (AliasCmd_Create) isAliasCmd() {}

type AliasCmd_Delete struct {
	Value DeleteCmd
}
func (AliasCmd_Delete) isAliasCmd() {}

type AliasCmd_Debug struct {
	Value DebugCmd
}
func (AliasCmd_Debug) isAliasCmd() {}

func AliasCmdCreate(value CreateCmd) AliasCmd {
	return AliasCmd_Create{Value: value}
}
func AliasCmdDelete(value DeleteCmd) AliasCmd {
	return AliasCmd_Delete{Value: value}
}
func AliasCmdDebug(value DebugCmd) AliasCmd {
	return AliasCmd_Debug{Value: value}
}

func TestParseBasicArg(t *testing.T) {
	sub, _err0 := ParseArgs[SingleCmd]([]string{"Alice"})
	if _err0 != nil {
		err := _err0
		t.Fatalf("unexpected error: %v", err)
	}

	switch __v := sub.(type) {
	case SingleCmd_Greet:
		cmd := __v.Value
		_ = cmd
		if cmd.Name != "Alice" {
			t.Errorf("expected Name='Alice', got '%s'", cmd.Name)
		}
	}
}

func TestParseFlag(t *testing.T) {
	sub, _err0 := ParseArgs[SingleCmd]([]string{"-c", "5", "Bob"})
	if _err0 != nil {
		err := _err0
		t.Fatalf("unexpected error: %v", err)
	}

	switch __v := sub.(type) {
	case SingleCmd_Greet:
		cmd := __v.Value
		_ = cmd
		if cmd.Count != 5 {
			t.Errorf("expected Count=5, got %d", cmd.Count)
		}
		if cmd.Name != "Bob" {
			t.Errorf("expected Name='Bob', got '%s'", cmd.Name)
		}
	}
}

func TestParseLongFlag(t *testing.T) {
	sub, _err0 := ParseArgs[SingleCmd]([]string{"--count", "3", "--verbose", "Charlie"})
	if _err0 != nil {
		err := _err0
		t.Fatalf("unexpected error: %v", err)
	}

	switch __v := sub.(type) {
	case SingleCmd_Greet:
		cmd := __v.Value
		_ = cmd
		if cmd.Count != 3 {
			t.Errorf("expected Count=3, got %d", cmd.Count)
		}
		if (!cmd.Verbose) {
			t.Error("expected Verbose=true")
		}
	}
}

func TestParseFlagWithEquals(t *testing.T) {
	sub, _err0 := ParseArgs[SingleCmd]([]string{"--count=7", "Dave"})
	if _err0 != nil {
		err := _err0
		t.Fatalf("unexpected error: %v", err)
	}

	switch __v := sub.(type) {
	case SingleCmd_Greet:
		cmd := __v.Value
		_ = cmd
		if cmd.Count != 7 {
			t.Errorf("expected Count=7, got %d", cmd.Count)
		}
	}
}

func TestParseDefaultValue(t *testing.T) {
	sub, _err0 := ParseArgs[SingleCmd]([]string{"Eve"})
	if _err0 != nil {
		err := _err0
		t.Fatalf("unexpected error: %v", err)
	}

	switch __v := sub.(type) {
	case SingleCmd_Greet:
		cmd := __v.Value
		_ = cmd
		if cmd.Count != 1 {
			t.Errorf("expected default Count=1, got %d", cmd.Count)
		}
	}
}

func TestParseSubcommandAdd(t *testing.T) {
	sub, _err0 := ParseArgs[MultiCmd]([]string{"add", "milk"})
	if _err0 != nil {
		err := _err0
		t.Fatalf("unexpected error: %v", err)
	}

	switch __v := sub.(type) {
	case MultiCmd_Add:
		cmd := __v.Value
		_ = cmd
		if cmd.Item != "milk" {
			t.Errorf("expected Item='milk', got '%s'", cmd.Item)
		}
	case MultiCmd_List:
		cmd := __v.Value
		_ = cmd
		t.Error("expected Add variant, got List")
	}
}

func TestParseSubcommandList(t *testing.T) {
	sub, _err0 := ParseArgs[MultiCmd]([]string{"list", "--all"})
	if _err0 != nil {
		err := _err0
		t.Fatalf("unexpected error: %v", err)
	}

	switch __v := sub.(type) {
	case MultiCmd_Add:
		cmd := __v.Value
		_ = cmd
		t.Error("expected List variant, got Add")
	case MultiCmd_List:
		cmd := __v.Value
		_ = cmd
		if (!cmd.All) {
			t.Error("expected All=true")
		}
	}
}

func TestParseSubWithParent(t *testing.T) {
	parent, sub, _err0 := ParseArgsSub[ParentCmd, MultiCmd]([]string{"-v", "add", "eggs"})
	if _err0 != nil {
		err := _err0
		t.Fatalf("unexpected error: %v", err)
	}

	if (!parent.Verbose) {
		t.Error("expected parent Verbose=true")
	}

	switch __v := sub.(type) {
	case MultiCmd_Add:
		cmd := __v.Value
		_ = cmd
		if cmd.Item != "eggs" {
			t.Errorf("expected Item='eggs', got '%s'", cmd.Item)
		}
	case MultiCmd_List:
		cmd := __v.Value
		_ = cmd
		t.Error("expected Add variant, got List")
	}
}

func TestParseMissingRequired(t *testing.T) {
	_, err := ParseArgs[SingleCmd]([]string{})
	if err == nil {
		t.Error("expected error for missing required argument")
	}

	parseErr, ok := err.(ParseError)
	if (!ok) {
		t.Fatalf("expected ParseError, got %T", err)
	}

	switch parseErr.Kind.(type) {
	case ErrorKind_MissingRequired:
	default:
		// Expected
		t.Errorf("expected MissingRequired error, got %v", parseErr.Kind)
	}
}

func TestParseUnknownFlag(t *testing.T) {
	_, err := ParseArgs[SingleCmd]([]string{"--unknown", "Frank"})
	if err == nil {
		t.Error("expected error for unknown flag")
	}

	parseErr, ok := err.(ParseError)
	if (!ok) {
		t.Fatalf("expected ParseError, got %T", err)
	}

	switch parseErr.Kind.(type) {
	case ErrorKind_UnknownFlag:
	default:
		// Expected
		t.Errorf("expected UnknownFlag error, got %v", parseErr.Kind)
	}
}

func TestParseMissingSubcommand(t *testing.T) {
	_, err := ParseArgs[MultiCmd]([]string{})
	if err == nil {
		t.Error("expected error for missing subcommand")
	}

	parseErr, ok := err.(ParseError)
	if (!ok) {
		t.Fatalf("expected ParseError, got %T", err)
	}

	switch parseErr.Kind.(type) {
	case ErrorKind_UnknownSubcommand:
	default:
		// Expected
		t.Errorf("expected UnknownSubcommand error, got %v", parseErr.Kind)
	}
}

func TestParseHelpFlag(t *testing.T) {
	_, err := ParseArgs[SingleCmd]([]string{"--help"})
	if err == nil {
		t.Error("expected error for help flag")
	}

	parseErr, ok := err.(ParseError)
	if (!ok) {
		t.Fatalf("expected ParseError, got %T", err)
	}

	switch parseErr.Kind.(type) {
	case ErrorKind_HelpRequested:
	default:
		// Expected
		t.Errorf("expected HelpRequested error, got %v", parseErr.Kind)
	}
}

func TestParseInvalidIntValue(t *testing.T) {
	_, err := ParseArgs[SingleCmd]([]string{"--count", "notanumber", "Grace"})
	if err == nil {
		t.Error("expected error for invalid int value")
	}

	parseErr, ok := err.(ParseError)
	if (!ok) {
		t.Fatalf("expected ParseError, got %T", err)
	}

	switch parseErr.Kind.(type) {
	case ErrorKind_InvalidValue:
	default:
		// Expected
		t.Errorf("expected InvalidValue error, got %v", parseErr.Kind)
	}
}

// Test parsing with alias
func TestParseAlias(t *testing.T) {
	// "rm" is an alias for "delete"
	sub, _err0 := ParseArgs[AliasCmd]([]string{"rm", "foo"})
	if _err0 != nil {
		err := _err0
		t.Fatalf("unexpected error: %v", err)
	}

	switch __v := sub.(type) {
	case AliasCmd_Delete:
		cmd := __v.Value
		_ = cmd
		if cmd.Name != "foo" {
			t.Errorf("expected Name='foo', got '%s'", cmd.Name)
		}
	default:
		t.Error("expected Delete variant")
	}
}

// Test parsing with second alias
func TestParseSecondAlias(t *testing.T) {
	// "del" is also an alias for "delete"
	sub, _err0 := ParseArgs[AliasCmd]([]string{"del", "bar"})
	if _err0 != nil {
		err := _err0
		t.Fatalf("unexpected error: %v", err)
	}

	switch __v := sub.(type) {
	case AliasCmd_Delete:
		cmd := __v.Value
		_ = cmd
		if cmd.Name != "bar" {
			t.Errorf("expected Name='bar', got '%s'", cmd.Name)
		}
	default:
		t.Error("expected Delete variant")
	}
}

// Test parsing with canonical name still works
func TestParseCanonicalName(t *testing.T) {
	sub, _err0 := ParseArgs[AliasCmd]([]string{"delete", "baz"})
	if _err0 != nil {
		err := _err0
		t.Fatalf("unexpected error: %v", err)
	}

	switch __v := sub.(type) {
	case AliasCmd_Delete:
		cmd := __v.Value
		_ = cmd
		if cmd.Name != "baz" {
			t.Errorf("expected Name='baz', got '%s'", cmd.Name)
		}
	default:
		t.Error("expected Delete variant")
	}
}

// Test hidden subcommand is still usable
func TestParseHiddenSubcommand(t *testing.T) {
	sub, _err0 := ParseArgs[AliasCmd]([]string{"debug"})
	if _err0 != nil {
		err := _err0
		t.Fatalf("unexpected error: %v", err)
	}

	switch __v := sub.(type) {
	case AliasCmd_Debug:
		cmd := __v.Value
		_ = cmd
		_ = cmd // Debug command parsed successfully
	default:
		t.Error("expected Debug variant")
	}
}

// Test hidden subcommand not in missing subcommand error
func TestHiddenNotInError(t *testing.T) {
	_, err := ParseArgs[AliasCmd]([]string{})
	if err == nil {
		t.Error("expected error for missing subcommand")
	}

	parseErr, ok := err.(ParseError)
	if (!ok) {
		t.Fatalf("expected ParseError, got %T", err)
	}

	// Error message should not contain "debug"
	if strings.Contains(parseErr.Message, "debug") {
		t.Errorf("hidden subcommand 'debug' should not appear in error message: %s", parseErr.Message)
	}

	// But should contain visible commands
	if (!strings.Contains(parseErr.Message, "create")) {
		t.Errorf("expected 'create' in error message: %s", parseErr.Message)
	}
	if (!strings.Contains(parseErr.Message, "delete")) {
		t.Errorf("expected 'delete' in error message: %s", parseErr.Message)
	}
}

func init() {
	runtime.RegisterAttr("slap.GreetTestCmd", "", Command{Name: "greet", About: "Greet command"})
	runtime.RegisterAttr("slap.GreetTestCmd", "Name", Arg{Position: 0, Help: "Name argument"})
	runtime.RegisterAttr("slap.GreetTestCmd", "Count", Flag{Short: "c", Long: "count", Default: "1"})
	runtime.RegisterAttr("slap.GreetTestCmd", "Verbose", Flag{Short: "v", Long: "verbose"})
	runtime.RegisterAttr("slap.SingleCmd", "", Subcommands{})
	runtime.RegisterAttr("slap.AddItemCmd", "", Command{Name: "add", About: "Add item"})
	runtime.RegisterAttr("slap.AddItemCmd", "Item", Arg{Position: 0})
	runtime.RegisterAttr("slap.ListItemsCmd", "", Command{Name: "list", About: "List items"})
	runtime.RegisterAttr("slap.ListItemsCmd", "All", Flag{Short: "a", Long: "all"})
	runtime.RegisterAttr("slap.MultiCmd", "", Subcommands{})
	runtime.RegisterAttr("slap.ParentCmd", "", Command{Name: "app"})
	runtime.RegisterAttr("slap.ParentCmd", "Verbose", Flag{Short: "v", Long: "verbose", Global: true})
	runtime.RegisterAttr("slap.CreateCmd", "", Command{Name: "create", About: "Create something"})
	runtime.RegisterAttr("slap.CreateCmd", "Name", Arg{Position: 0})
	runtime.RegisterAttr("slap.DeleteCmd", "", Command{Name: "delete", About: "Delete something"})
	runtime.RegisterAttr("slap.DeleteCmd", "Name", Arg{Position: 0})
	runtime.RegisterAttr("slap.DebugCmd", "", Command{Name: "debug", About: "Debug mode"})
	runtime.RegisterAttr("slap.AliasCmd", "", Subcommands{})
	runtime.RegisterAttr("slap.AliasCmd", "Delete", Alias{Name: "rm"})
	runtime.RegisterAttr("slap.AliasCmd", "Delete", Alias{Name: "del"})
	runtime.RegisterAttr("slap.AliasCmd", "Debug", Hidden{})
	runtime.RegisterAttr("slap.SingleCmd", "Greet", runtime.EnumVariant{WrapperType: SingleCmd_Greet{}})
	runtime.RegisterAttr("slap.MultiCmd", "Add", runtime.EnumVariant{WrapperType: MultiCmd_Add{}})
	runtime.RegisterAttr("slap.MultiCmd", "List", runtime.EnumVariant{WrapperType: MultiCmd_List{}})
	runtime.RegisterAttr("slap.AliasCmd", "Create", runtime.EnumVariant{WrapperType: AliasCmd_Create{}})
	runtime.RegisterAttr("slap.AliasCmd", "Delete", runtime.EnumVariant{WrapperType: AliasCmd_Delete{}})
	runtime.RegisterAttr("slap.AliasCmd", "Debug", runtime.EnumVariant{WrapperType: AliasCmd_Debug{}})
}
