package clipboard

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/atotto/clipboard"
)

// PasteText pastes the given text using clipboard and simulates Cmd+V or Ctrl+V
func PasteText(text string) error {
	fmt.Println("Pasting text at", time.Now().Format(time.RFC3339), ":", text)

	if err := clipboard.WriteAll(text); err != nil {
		return fmt.Errorf("failed to write to clipboard: %w", err)
	}

	// Add a small delay to ensure clipboard is ready
	time.Sleep(100 * time.Millisecond)

	return simulatePaste()
}

// FocusPasteText switches focus away and back, then pastes text.
func FocusPasteText(text string) error {
	fmt.Println("Hide-paste triggered at", time.Now().Format(time.RFC3339), ":", text)

	if err := clipboard.WriteAll(text); err != nil {
		return fmt.Errorf("failed to write to clipboard: %w", err)
	}

	// Add a small delay to ensure clipboard is ready
	time.Sleep(100 * time.Millisecond)

	switch runtime.GOOS {
	case "darwin":
		cmd := exec.Command("osascript", "-e", `
			tell application "System Events"
				-- Switch to next application
				keystroke tab using command down
				delay 0.3
				-- Paste in the new application
				keystroke "v" using command down
			end tell`)
		cmd.Stderr = os.Stderr
		return cmd.Run()

	case "windows":
		cmd := exec.Command("powershell", "-Command", `
			Add-Type -AssemblyName System.Windows.Forms;
			[System.Windows.Forms.SendKeys]::SendWait('%%{TAB}');
			Start-Sleep -Milliseconds 300;
			[System.Windows.Forms.SendKeys]::SendWait('^v');
		`)
		cmd.Stderr = os.Stderr
		return cmd.Run()

	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// simulatePaste sends a paste command based on platform.
func simulatePaste() error {
	switch runtime.GOOS {
	case "darwin":
		cmd := exec.Command("osascript", "-e", `
			tell application "System Events"
				keystroke "v" using command down
			end tell`)
		cmd.Stderr = os.Stderr
		return cmd.Run()

	case "windows":
		cmd := exec.Command("powershell", "-Command", `
			Add-Type -AssemblyName System.Windows.Forms;
			[System.Windows.Forms.SendKeys]::SendWait('^v');
		`)
		cmd.Stderr = os.Stderr
		return cmd.Run()

	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}
