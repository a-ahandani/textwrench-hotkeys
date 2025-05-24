//go:build windows

package clipboard

import (
	"errors"
	"fmt"
	"time"

	"github.com/atotto/clipboard"
	"github.com/micmonay/keybd_event"
)

// sendKeyCombo presses and releases Ctrl+<key> to trigger copy/paste operations.
func sendKeyCombo(key int) error {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return fmt.Errorf("failed to create key bonding: %w", err)
	}
	// Set the key to send
	kb.SetKeys(key)
	// Indicate that Ctrl should be held
	kb.HasCTRL(true)
	// A small delay between press and release
	time.Sleep(100 * time.Millisecond)

	// Launching will press Ctrl, the key, then release both
	if err := kb.Launching(); err != nil {
		return fmt.Errorf("failed to send key combo: %w", err)
	}
	return nil
}

// readSelectedText clears the clipboard, sends Ctrl+C, and returns the copied text.
func readSelectedText() (string, error) {
	// Clear the clipboard
	if err := clipboard.WriteAll(""); err != nil {
		return "", fmt.Errorf("failed to clear clipboard: %w", err)
	}

	// Allow time for clipboard to clear
	time.Sleep(100 * time.Millisecond)

	// Send Ctrl+C
	if err := sendKeyCombo(keybd_event.VK_C); err != nil {
		return "", fmt.Errorf("failed to send Ctrl+C: %w", err)
	}

	// Poll for up to 500ms until we get non-empty clipboard content
	var text string
	start := time.Now()
	for time.Since(start) < 500*time.Millisecond {
		text, _ = clipboard.ReadAll()
		if text != "" {
			fmt.Println("📋 Selected:", text)
			return text, nil
		}
		// Brief pause before retrying
		time.Sleep(50 * time.Millisecond)
	}

	// If still empty, return an error
	return "", errors.New("no text copied from selection")
}

// writeText sets the clipboard to the provided text and sends Ctrl+V to paste.
func writeText(text string) error {
	// Update clipboard content directly
	if err := clipboard.WriteAll(text); err != nil {
		return fmt.Errorf("failed to set clipboard: %w", err)
	}

	// Allow time for clipboard update
	time.Sleep(100 * time.Millisecond)

	// Send Ctrl+V to paste
	if err := sendKeyCombo(keybd_event.VK_V); err != nil {
		return fmt.Errorf("failed to send Ctrl+V: %w", err)
	}

	fmt.Println("📋 Pasted:", text)
	return nil
}
