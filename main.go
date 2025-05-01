//go:build windows

package main

import (
	"fmt"
	"log"
	"time"
	"unsafe"

	"github.com/atotto/clipboard"
	"github.com/micmonay/keybd_event"
	"golang.org/x/sys/windows"
)

const (
	modCtrl   = 0x0002
	modShift  = 0x0004
	vkC       = 0x43
	idHotkey  = 1
	WM_HOTKEY = 0x0312
)

var (
	user32         = windows.NewLazySystemDLL("user32.dll")
	registerHotKey = user32.NewProc("RegisterHotKey")
	getMessage     = user32.NewProc("GetMessageW")
)

type msg struct {
	Hwnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct{ X, Y int32 }
}

func main() {
	ret, _, err := registerHotKey.Call(0, idHotkey, modCtrl|modShift, vkC)
	if ret == 0 {
		log.Fatalf("Failed to register hotkey: %v", err)
	}
	fmt.Println("✅ Hotkey registered: Ctrl+Shift+C")

	var m msg
	for {
		getMessage.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if m.Message == WM_HOTKEY && m.WParam == idHotkey {
			handleHotkey()
		}
	}
}

func handleHotkey() {
	fmt.Println("🔁 Hotkey triggered")

	// Clear clipboard to avoid stale data
	clipboard.WriteAll("")

	// Simulate Ctrl+C to get selected text
	if err := sendCtrlKey(keybd_event.VK_C); err != nil {
		log.Println("❌ Failed to send Ctrl+C:", err)
		return
	}

	time.Sleep(70 * time.Millisecond) // wait for clipboard to update

	selected, err := clipboard.ReadAll()
	if err != nil || selected == "" {
		log.Println("❌ Could not read clipboard or empty selection")
		return
	}

	fmt.Println("📋 Selected:", selected)

	processed := selected + " [processed]"

	// Copy processed text to clipboard
	if err := clipboard.WriteAll(processed); err != nil {
		log.Println("❌ Failed to write to clipboard:", err)
		return
	}

	time.Sleep(100 * time.Millisecond)

	// Simulate Ctrl+V to paste replacement
	if err := sendCtrlKey(keybd_event.VK_V); err != nil {
		log.Println("❌ Failed to send Ctrl+V:", err)
	}
}

func sendCtrlKey(key int) error {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return err
	}
	kb.SetKeys(key)
	kb.HasCTRL(true)
	return kb.Launching()
}
