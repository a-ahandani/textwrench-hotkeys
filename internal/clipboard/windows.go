//go:build windows

package clipboard

import (
	"fmt"
	"time"

	"github.com/micmonay/keybd_event"
)

func sendKeyCombo(key int) error {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return fmt.Errorf("failed to create key bonding: %w", err)
	}

	kb.SetKeys(key)
	kb.HasCTRL(true)

	time.Sleep(100 * time.Millisecond)

	if err := kb.Launching(); err != nil {
		return fmt.Errorf("failed to send key combo: %w", err)
	}
	return nil
}

func switchToNextApp() error {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return fmt.Errorf("failed to create key bonding: %w", err)
	}

	kb.SetKeys(keybd_event.VK_TAB)
	kb.HasALT(true)

	time.Sleep(100 * time.Millisecond)

	if err := kb.Launching(); err != nil {
		return fmt.Errorf("failed to switch to next app: %w", err)
	}
	return nil
}
