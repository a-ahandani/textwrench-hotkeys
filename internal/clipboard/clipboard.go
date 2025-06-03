package clipboard

import (
	"errors"
	"fmt"
	"time"

	"github.com/atotto/clipboard"
	"github.com/micmonay/keybd_event"
)

func ReadSelectedText() (string, error) {
	if err := clipboard.WriteAll(""); err != nil {
		return "", fmt.Errorf("failed to clear clipboard: %w", err)
	}
	time.Sleep(100 * time.Millisecond)

	if err := sendKeyCombo(keybd_event.VK_C); err != nil {
		return "", fmt.Errorf("failed to send Super+C: %w", err)
	}

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

func WriteText(text string) error {
	if err := clipboard.WriteAll(text); err != nil {
		return fmt.Errorf("failed to set clipboard: %w", err)
	}
	time.Sleep(100 * time.Millisecond)

	if err := sendKeyCombo(keybd_event.VK_V); err != nil {
		return fmt.Errorf("failed to send Ctrl+V: %w", err)
	}

	fmt.Println("📋 Pasted:", text)
	return nil
}

func FocusPasteText(text string) error {
	fmt.Println("Hide-paste triggered at", time.Now().Format(time.RFC3339), ":", text)

	if err := clipboard.WriteAll(text); err != nil {
		return fmt.Errorf("failed to write to clipboard: %w", err)
	}
	time.Sleep(100 * time.Millisecond)

	return switchToNextApp()
}
