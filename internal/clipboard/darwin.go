//go:build darwin

package clipboard

import (
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/atotto/clipboard"
	"github.com/micmonay/keybd_event"
)

func sendKeyCombo(key int) error {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return fmt.Errorf("failed to create key bonding: %w", err)
	}

	kb.SetKeys(key)

	switch runtime.GOOS {
	case "darwin":
		kb.HasSuper(true)
	case "windows":
		kb.HasCTRL(true)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	time.Sleep(100 * time.Millisecond)

	if err := kb.Launching(); err != nil {
		return fmt.Errorf("failed to send key combo: %w", err)
	}
	return nil
}

func readSelectedText() (string, error) {
	// Clear the clipboard
	if err := clipboard.WriteAll(""); err != nil {
		return "", fmt.Errorf("failed to clear clipboard: %w", err)
	}

	// Allow time for clipboard to clear
	time.Sleep(100 * time.Millisecond)

	if err := sendKeyCombo(keybd_event.VK_C); err != nil {
		return "", fmt.Errorf("failed to send Super+C: %w", err)
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
		time.Sleep(50 * time.Millisecond)
	}

	return "", errors.New("no text copied from selection")
}

func writeText(text string) error {
	if err := clipboard.WriteAll(text); err != nil {
		return fmt.Errorf("failed to set clipboard: %w", err)
	}

	time.Sleep(100 * time.Millisecond)

	// Send Ctrl+V to paste
	if err := sendKeyCombo(keybd_event.VK_V); err != nil {
		return fmt.Errorf("failed to send Ctrl+V: %w", err)
	}

	fmt.Println("📋 Pasted:", text)
	return nil
}
