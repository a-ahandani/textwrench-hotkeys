package clipboard

import (
	"fmt"
	"time"
)

func ReadSelectedText() (string, error) {
	fmt.Println("Reading selected text at", time.Now().Format(time.RFC3339))
	return readSelectedText()
}

func WriteText(text string) error {
	return writeText(text)
}
