package readline

// viinsKeys are the default keymaps in Vim Insert mode
var viinsKeys = keymap{
	// Standard
	"^[":      "vi-cmd-mode",
	"^M":      "accept-line",
	"^L":      "clear-screen",
	"^Y":      "yank",
	"^E":      "end-of-line",
	"^A":      "beginning-of-line",
	"^B":      "backward-char",
	"^F":      "forward-char",
	"^K":      "kill-line",
	"^N":      "down-line-or-history",
	"^P":      "up-line-or-history",
	"^W":      "backward-kill-word",
	"^?":      "backward-delete-char",
	"^_":      "undo",
	"^[[3~":   "delete-char",
	"^[[H":    "beginning-of-line",
	"^[[F":    "end-of-line",
	"^[[5~":   "history-search", // TODO
	"^[[6~":   "menu-select",    // TODO
	"^[[C":    "vi-forward-char",
	"^[[D":    "vi-backward-char",
	"^[[3;5~": "kill-word",
	"^[[1;5C": "forward-word",
	"^[[1;5D": "backward-word",
	"^[[A":    "up-line-or-search",   // TODO
	"^[[B":    "down-line-or-select", // TODO
	" ":       "space",
	"[!-~]":   "self-insert", // Any non-empty, non-modified key (no control sequences)
	// `[^\-^^]`: "self-insert",
}

// viinsKeymaps are the default keymaps in Vim Command mode
var vicmdKeys = keymap{
	// Standard
	"i":       "vi-insert-mode",
	"^M":      "accept-line",
	"^L":      "clear-screen",
	"^N":      "down-history",
	"^P":      "up-history",
	"^[[3~":   "delete-char",
	"^[[6~":   "down-line-or-history",
	"^[[5~":   "up-line-or-history",
	"^[[H":    "beginning-of-line",
	"^[[F":    "end-of-line",
	"^[[A":    "history-search",
	"^[[B":    "menu-select",
	"^[[3;5~": "kill-word",
	"^[[1;5C": "forward-word",
	"^[[1;5D": "backward-word",

	// History

	// Vim
	"^A":    "switch-keyword",
	"^X":    "switch-keyword",
	"^R":    "redo",
	"^?":    "backward-delete-char",
	"^[[C":  "vi-forward-char",
	"^[[D":  "vi-backward-char",
	" ":     "vi-forward-char",
	"$":     "vi-end-of-line",
	"%":     "vi-match-bracket",
	"\"":    "vi-set-buffer",
	"0":     "vi-digit-or-beginning-of-line",
	"a":     "vi-add-next",
	"A":     "vi-add-eol",
	"b":     "vi-backward-word",
	"B":     "vi-backward-blank-word",
	"c":     "vi-change",
	"C":     "vi-change-eol",
	"d":     "vi-delete",
	"D":     "vi-kill-eol",
	"e":     "vi-forward-word-end",
	"E":     "vi-forward-blank-word-end",
	"f":     "vi-find-next-char",
	"t":     "vi-find-next-char-skip",
	"I":     "vi-insert-bol",
	"h":     "vi-backward-char",
	"l":     "vi-forward-char",
	"j":     "down-line-or-history",
	"k":     "up-line-or-history",
	"p":     "vi-put-after",
	"P":     "vi-put-before",
	"r":     "vi-replace-chars",
	"R":     "vi-replace",
	"F":     "vi-find-prev-char",
	"T":     "vi-find-prev-char-skip",
	"s":     "vi-substitute",
	"u":     "undo",
	"v":     "visual-mode",
	"V":     "visual-line-mode",
	"w":     "vi-forward-word",
	"W":     "vi-forward-blank-word",
	"x":     "vi-delete-char",
	"X":     "vi-backward-delete-char",
	"y":     "vi-yank",
	"Y":     "vi-yank-whole-line",
	"|":     "vi-goto-column",
	"~":     "vi-swap-case",
	"g~":    "vi-oper-swap-case",
	`[1-9]`: "digit-argument",
}

// viinsKeymaps are the default keymaps in Vim Operating Pending mode
var vioppKeys = keymap{
	"^[": "vi-cmd-mode",
	"aW": "select-a-blank-word",
	"aa": "select-a-shell-word",
	"aw": "select-a-word",
	"iW": "select-in-blank-word",
	"ia": "select-in-shell-word",
	"iw": "select-in-word",
	"j":  "down-line", // Not sure since no multiline
	"k":  "up-line",   // Not sure since no multiline
}

// viinsKeymaps are the default keymaps in Vim Visual mode
var visualKeys = keymap{
	"^[": "vi-cmd-mode",
	"aW": "select-a-blank-word",
	"aa": "select-a-shell-word",
	"aw": "select-a-word",
	"iW": "select-in-blank-word",
	"ia": "select-in-shell-word",
	"iw": "select-in-word",
	"s":  "vi-substitute",
	"S":  "vi-add-surround",
	"a":  "vi-select-surround",
	"c":  "vi-change",
	"d":  "vi-delete",
	"i":  "vi-select-surround",
	"j":  "down-line", // Not sure since no multiline
	"k":  "up-line",   // Not sure since no multiline
	"u":  "vi-down-case",
	"v":  "vi-edit-command-line",
	"x":  "vi-delete",
	"y":  "vi-yank",
	"~":  "vi-swap-case",
}

// changeMovements is used for some widgets that only
// accept movement widgets as arguments (like vi-change)
var changeMovements = map[string]string{
	"$": "vi-end-of-line",
	"%": "vi-match-bracket",
	"^": "vi-first-non-blank",
	"0": "beginning-of-line",
	"b": "vi-backward-word",
	"B": "vi-backward-blank-word",
	"w": "vi-forward-word",
	"W": "vi-forward-blank-word",
	"e": "vi-forward-word-end",
	"E": "vi-forward-blank-word-end",
	"f": "vi-find-next-char",
	"F": "vi-find-prev-char",
	"t": "vi-find-next-char-skip",
	"T": "vi-find-prev-char-skip",
	"g": "",
	"h": "vi-backward-char",
	"l": "vi-forward-char",
	"s": "vi-change-surround",
	"a": "vi-select-surround",
	"i": "vi-select-surround",
}
