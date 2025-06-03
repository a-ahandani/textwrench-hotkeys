//go:build windows

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
	kb.HasCTRL(true)

	time.Sleep(100 * time.Millisecond)

	if err := kb.Launching(); err != nil {
		return fmt.Errorf("failed to send key combo: %w", err)
	}
	return nil
}

func switchToNextApp() error {
	cmd := exec.Command("powershell", "-NoProfile", "-Command", `
Add-Type -TypeDefinition @"
using System;
using System.Runtime.InteropServices;
public class Keyboard {
    [DllImport("user32.dll")]
    public static extern void keybd_event(byte bVk, byte bScan, uint dwFlags, UIntPtr dwExtraInfo);
    const int KEYEVENTF_KEYUP = 0x0002;
    const byte VK_MENU = 0x12;
    const byte VK_TAB = 0x09;

    public static void SendAltTab() {
        keybd_event(VK_MENU, 0, 0, UIntPtr.Zero);
        keybd_event(VK_TAB, 0, 0, UIntPtr.Zero);
        keybd_event(VK_TAB, 0, KEYEVENTF_KEYUP, UIntPtr.Zero);
        keybd_event(VK_MENU, 0, KEYEVENTF_KEYUP, UIntPtr.Zero);
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
}
