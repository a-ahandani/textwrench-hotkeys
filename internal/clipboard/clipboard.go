package clipboard

// ReadSelectedText triggers Cmd/Ctrl+C and returns selected clipboard text.
func ReadSelectedText() (string, error) {
	return readSelectedText()
}

// WriteText writes the given text to the clipboard and triggers paste (Cmd/Ctrl+V).
func WriteText(text string) error {
	return writeText(text)
}
