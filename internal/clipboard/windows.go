//go:build windows

package clipboard

import (
	"time"

	"github.com/atotto/clipboard"
	"github.com/micmonay/keybd_event"
)

func readSelectedText() (string, error) {
	if err := sendCopy(); err != nil {
		return "", err
	}
	time.Sleep(100 * time.Millisecond)
	return clipboard.ReadAll()
}

func writeText(text string) error {
	if err := clipboard.WriteAll(text); err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond)
	return sendPaste()
}

func sendCopy() error {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return err
	}
	kb.HasCTRL(true)
	kb.SetKeys(keybd_event.VK_C)
	return kb.Launching()
}

func sendPaste() error {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return err
	}
	kb.HasCTRL(true)
	kb.SetKeys(keybd_event.VK_V)
	return kb.Launching()
}
