package readline

import (
	"fmt"

	"github.com/reeflective/readline/inputrc"
	"github.com/reeflective/readline/internal/completion"
	"github.com/reeflective/readline/internal/core"
	"github.com/reeflective/readline/internal/display"
	"github.com/reeflective/readline/internal/editor"
	"github.com/reeflective/readline/internal/history"
	"github.com/reeflective/readline/internal/keymap"
	"github.com/reeflective/readline/internal/macro"
	"github.com/reeflective/readline/internal/term"
	"github.com/reeflective/readline/internal/ui"
)

// Shell is used to encapsulate the parameter group and run time of any given
// readline instance so that you can reuse the readline API for multiple entry
// captures without having to repeatedly unload configuration.
type Shell struct {
	// Core editor
	line       *core.Line       // The input line buffer and its management methods.
	cursor     *core.Cursor     // The cursor and its medhods.
	selection  *core.Selection  // The selection managees various visual/pending selections.
	Iterations *core.Iterations // Digit arguments for repeating commands.
	Buffers    *editor.Buffers  // buffers (Vim registers) and methods use/manage/query them.
	Keys       *core.Keys       // Keys is in charge of reading and managing buffered user input.
	Keymap     *keymap.Engine   // Manages main/local keymaps, binds, stores command functions, etc.

	// User interface
	Config      *inputrc.Config    // Contains all keymaps, binds and per-application settings.
	Opts        []inputrc.Option   // Inputrc file parsing options (app/term/values, etc).
	Prompt      *ui.Prompt         // The prompt engine computes and renders prompt strings.
	Hint        *ui.Hint           // Usage/hints for completion/isearch below the input line.
	Completions *completion.Engine // Completions generation and display.
	History     *history.Sources   // History manages all history types/sources (past commands and undo)
	Macros      *macro.Engine      // Record, use and display macros.
	Display     *display.Engine    // Manages display refresh/update/clearing.

	// User-provided functions

	// AcceptMultiline enables the caller to decide if the shell should keep reading
	// for user input on a new line (therefore, with the secondary prompt), or if it
	// should return the current line at the end of the `rl.Readline()` call.
	// This function should return true if the line is deemed complete (thus asking
	// the shell to return from its Readline() loop), or false if the shell should
	// keep reading input on a newline (thus, insert a newline and read).
	AcceptMultiline func(line []rune) (accept bool)

	// SyntaxHighlighter is a helper function to provide syntax highlighting.
	// Once enabled, set to nil to disable again.
	SyntaxHighlighter func(line []rune) string

	// Completer is a function that produces completions.
	// It takes the readline line ([]rune) and cursor pos as parameters,
	// and returns completions with their associated metadata/settings.
	Completer func(line []rune, cursor int) Completions

	// SyntaxCompletion is used to autocomplete code syntax (like braces and
	// quotation marks). If you want to complete words or phrases then you might
	// be better off using the Completer function.
	// SyntaxCompletion takes the line ([]rune) and cursor position, and returns
	// the new line and cursor position.
	SyntaxCompleter func(line []rune, cursor int) ([]rune, int)

	// HintText is a helper function which displays hint text below the line.
	// HintText takes the line input from the promt and the cursor position.
	// It returns the hint text to display.
	HintText func(line []rune, cursor int) []rune
}

// NewShell returns a readline shell instance initialized with a default
// inputrc configuration and binds, and with an in-memory command history.
// The constructor accepts an optional list of inputrc configuration options,
// which are used when parsing/loading and applying any inputrc configuration.
func NewShell(opts ...inputrc.Option) *Shell {
	shell := new(Shell)

	// Core editor
	keys := new(core.Keys)
	line := new(core.Line)
	cursor := core.NewCursor(line)
	selection := core.NewSelection(line, cursor)
	iterations := new(core.Iterations)

	shell.Keys = keys
	shell.line = line
	shell.cursor = cursor
	shell.selection = selection
	shell.Buffers = editor.NewBuffers()
	shell.Iterations = iterations

	// Keymaps and commands
	keymaps, config := keymap.NewEngine(keys, iterations, opts...)
	keymaps.Register(shell.standardCommands())
	keymaps.Register(shell.viCommands())
	keymaps.Register(shell.historyCommands())
	keymaps.Register(shell.completionCommands())

	shell.Keymap = keymaps
	shell.Config = config
	shell.Opts = opts

	// User interface
	hint := new(ui.Hint)
	prompt := ui.NewPrompt(line, cursor, keymaps, config)
	macros := macro.NewEngine(keys, hint)
	completer := completion.NewEngine(keys, line, cursor, selection, hint, keymaps, config)
	history := history.NewSources(line, cursor, hint, config)
	display := display.NewEngine(keys, selection, history, prompt, hint, completer, config)

	completer.SetAutocompleter(shell.commandCompletion)

	shell.Config = config
	shell.Hint = hint
	shell.Prompt = prompt
	shell.Completions = completer
	shell.Macros = macros
	shell.History = history
	shell.Display = display

	return shell
}

// Line is the shell input line buffer.
// Contains methods to search and modify its contents,
// split itself with tokenizers, and displaying itself.
//
// When the shell is in incremental-search mode, this line is the minibuffer.
func (rl *Shell) Line() *core.Line { return rl.line }

// Cursor is the cursor position in the current line buffer.
// Contains methods to set, move, describe and check itself.
//
// When the shell is in incremental-search mode, this cursor is the minibuffer's one.
func (rl *Shell) Cursor() *core.Cursor { return rl.cursor }

// Selection contains all regions of an input line that are currently selected/marked
// with either a begin and/or end position. The main selection is the visual one, used
// with the default cursor mark and position, and contains a list of additional surround
// selections used to change/select multiple parts of the line at once.
func (rl *Shell) Selection() *core.Selection { return rl.selection }

// Log prints a formatted string below the current line and redisplays the prompt
// and input line (and possibly completions/hints if active) below the logged string.
// A newline is added to the message so that the prompt is correctly refreshed below.
func (rl *Shell) Log(msg string, args ...any) (n int, err error) {
	// First go back to the last line of the input line,
	// and clear everything below (hints and completions).
	rl.Display.CursorBelowLine()
	term.MoveCursorBackwards(term.GetWidth())
	fmt.Print(term.ClearScreenBelow)

	// Skip a line, and print the formatted message.
	n, err = fmt.Printf(msg+"\n", args...)

	// Redisplay the prompt, input line and active helpers.
	rl.Prompt.PrimaryPrint()
	rl.Display.Refresh()

	return
}

// LogTransient prints a formatted string in place of the current prompt and input
// line, and then refreshes, or "pushes" the prompt/line below this printed message.
func (rl *Shell) LogTransient(msg string, args ...any) (n int, err error) {
	// First go back to the beginning of the line/prompt, and
	// clear everything below (prompt/line/hints/completions).
	if rl.Prompt.Refreshing() {
		term.MoveCursorUp(1)
	}

	rl.Display.CursorToLineStart()
	term.MoveCursorBackwards(term.GetWidth())

	term.MoveCursorUp(rl.Prompt.PrimaryUsed())
	fmt.Print(term.ClearScreenBelow)

	// Print the logged message.
	n, err = fmt.Printf(msg+"\n", args...)

	// Redisplay the prompt, input line and active helpers.
	rl.Prompt.PrimaryPrint()
	rl.Display.Refresh()

	return
}
