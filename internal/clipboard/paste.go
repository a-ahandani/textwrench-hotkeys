package clipboard

import (
	"os/exec"
	"runtime"

	"github.com/atotto/clipboard"
)

func PasteText(text string) error {
	// 1. Set clipboard contents
	if err := clipboard.WriteAll(text); err != nil {
		return err
	}

	// 2. Simulate paste (Cmd+V or Ctrl+V)
	switch runtime.GOOS {
	case "darwin":
		cmd := exec.Command("osascript", "-e", `tell application "System Events" to keystroke "v" using command down`)
		return cmd.Run()
	case "windows":
		cmd := exec.Command("powershell", "-Command", `[System.Windows.Forms.SendKeys]::SendWait('^v')`)
		return cmd.Run()
	default:
		return nil
	}
}
