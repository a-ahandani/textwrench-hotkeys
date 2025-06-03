//go:build darwin

package clipboard

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/micmonay/keybd_event"
)

func sendKeyCombo(key int) error {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return fmt.Errorf("failed to create key bonding: %w", err)
	}

	kb.SetKeys(key)
	kb.HasSuper(true)

	time.Sleep(100 * time.Millisecond)

	if err := kb.Launching(); err != nil {
		return fmt.Errorf("failed to send key combo: %w", err)
	}
	return nil
}

func switchToNextApp() error {
	cmd := exec.Command("osascript", "-e", `
		tell application "System Events"
			keystroke tab using command down
		end tell`)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
