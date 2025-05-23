//go:build windows

package clipboard

import (
	"fmt"
	"log"
	"time"

	"github.com/atotto/clipboard"
	"github.com/micmonay/keybd_event"
)

func sendCtrlKey(key int) error {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return err
	}
	kb.SetKeys(key)
	kb.HasCTRL(true)
	return kb.Launching()
}

func readSelectedText() (string, error) {
	clipboard.WriteAll("")
	time.Sleep(100 * time.Millisecond)
	if err := sendCtrlKey(keybd_event.VK_C); err != nil {
		return "", err
	}

	time.Sleep(100 * time.Millisecond)

	selected, err := clipboard.ReadAll()
	fmt.Println("📋 Selected:", selected)
	if err != nil || selected == "" {
		return "", err
	}
	return selected, nil
}

func writeText(text string) error {
	fmt.Println("📋 Writing to clipboard:", text)
	time.Sleep(100 * time.Millisecond)

	if err := sendCtrlKey(keybd_event.VK_V); err != nil {
		log.Println("❌ Failed to send Ctrl+V:", err)
	}
	return nil
}
