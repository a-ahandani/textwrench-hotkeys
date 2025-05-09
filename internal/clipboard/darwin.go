//go:build darwin

package clipboard

import (
	"os/exec"
	"time"

	"github.com/atotto/clipboard"
)

func readSelectedText() (string, error) {
	if err := exec.Command("osascript", "-e", `tell application "System Events" to keystroke "c" using {command down}`).Run(); err != nil {
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
	return exec.Command("osascript", "-e", `tell application "System Events" to keystroke "v" using {command down}`).Run()
}
