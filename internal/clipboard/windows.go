//go:build windows

package clipboard

import (
	"bytes"
	"os/exec"
	"strings"
	"time"
)

func sendCtrlKey(key string) error {
	// Simulate Ctrl+C or Ctrl+V using PowerShell + Windows Script Host
	script := `
		Add-Type -AssemblyName System.Windows.Forms
		[System.Windows.Forms.SendKeys]::SendWait('^` + key + `')
	`
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	return cmd.Run()
}

func readSelectedText() (string, error) {
	// Simulate Ctrl+C
	if err := sendCtrlKey("c"); err != nil {
		return "", err
	}
	time.Sleep(100 * time.Millisecond)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", "Get-Clipboard")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimRight(out.String(), "\r\n"), nil
}

func writeText(text string) error {
	// Set clipboard content
	script := `Set-Clipboard -Value @"
` + text + `
"@`
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	if err := cmd.Run(); err != nil {
		return err
	}

	time.Sleep(100 * time.Millisecond)

	// Simulate Ctrl+V
	return sendCtrlKey("v")
}
