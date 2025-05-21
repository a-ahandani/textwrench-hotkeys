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
		cmd := exec.Command("powershell", "-NoProfile", "-Command", `
Add-Type -TypeDefinition @"
using System;
using System.Runtime.InteropServices;
public class Keyboard {
    [DllImport("user32.dll")]
    public static extern void keybd_event(byte bVk, byte bScan, uint dwFlags, UIntPtr dwExtraInfo);
    const int KEYEVENTF_KEYUP = 0x0002;
    const byte VK_MENU = 0x12;  // Alt
    const byte VK_TAB = 0x09;   // Tab

    public static void SendAltTab() {
        keybd_event(VK_MENU, 0, 0, UIntPtr.Zero);             // Alt down
        keybd_event(VK_TAB, 0, 0, UIntPtr.Zero);              // Tab down
        keybd_event(VK_TAB, 0, KEYEVENTF_KEYUP, UIntPtr.Zero); // Tab up
        keybd_event(VK_MENU, 0, KEYEVENTF_KEYUP, UIntPtr.Zero); // Alt up
    }
}
"@;

[Keyboard]::SendAltTab();
Start-Sleep -Milliseconds 400;
Add-Type -AssemblyName System.Windows.Forms;
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
