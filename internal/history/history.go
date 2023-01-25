package history

var defaultSourceName = "default history"

// Source is an interface to allow you to write your own history logging tools.
// By default readline will just use the dummyLineHistory interface which only
// logs the history to memory ([]string to be precise).
type Source interface {
	// Append takes the line and returns an updated number of lines or an error
	Write(string) (int, error)

	// GetLine takes the historic line number and returns the line or an error
	GetLine(int) (string, error)

	// Len returns the number of history lines
	Len() int

	// Dump returns everything in readline. The return is an interface{} because
	// not all LineHistory implementations will want to structure the history in
	// the same way. And since Dump() is not actually used by the readline API
	// internally, this methods return can be structured in whichever way is most
	// convenient for your own applications (or even just create an empty
	// function which returns `nil` if you don't require Dump() either)
	Dump() interface{}
}

// memory is an in memory history.
// One such history is bound to the readline shell by default.
type memory struct {
	items []string
}

// Write to history.
func (h *memory) Write(s string) (int, error) {
	h.items = append(h.items, s)
	return len(h.items), nil
}

// GetLine returns a line from history.
func (h *memory) GetLine(i int) (string, error) {
	return h.items[i], nil
}

// Len returns the number of lines in history.
func (h *memory) Len() int {
	return len(h.items)
}

// Dump returns the entire history.
func (h *memory) Dump() interface{} {
	return h.items
}

// null is a null History interface for when you don't
// want line entries remembered eg password input.
type null struct{}

// Write to history.
func (h *null) Write(s string) (int, error) {
	return 0, nil
}

// GetLine returns a line from history.
func (h *null) GetLine(i int) (string, error) {
	return "", nil
}

// Len returns the number of lines in history.
func (h *null) Len() int {
	return 0
}

// Dump returns the entire history.
func (h *null) Dump() interface{} {
	return []string{}
}
