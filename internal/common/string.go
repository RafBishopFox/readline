package common

import (
	"regexp"
	"unicode/utf8"

	"github.com/reeflective/readline/internal/term"
)

// This file contains some utilities to get compute various stuff on strings.

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var re = regexp.MustCompile(ansi)

func stripAnsi(str string) string {
	return re.ReplaceAllString(str, "")
}

func realLength(s string) int {
	colorStripped := stripAnsi(s)
	return utf8.RuneCountInString(colorStripped)
}

// lineSpan computes the number of columns and lines that are needed for a given line.
func lineSpan(line []rune, idx, indent int) (x, y int) {
	termWidth := term.GetWidth()
	lineLen := realLength(string(line))
	lineLen += indent

	cursorY := lineLen / termWidth
	cursorX := lineLen % termWidth

	// The very first (unreal) line counts for nothing,
	// so by opposition all others count for one more.
	if idx == 0 {
		cursorY--
	}

	// Any excess wrap means a newline.
	if cursorX > 0 {
		cursorY++
	}

	// Empty lines are still considered a line.
	if cursorY == 0 && idx != 0 {
		cursorY++
	}

	return cursorX, cursorY
}
