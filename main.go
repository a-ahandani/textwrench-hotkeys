//go:build windows

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
	"unsafe"

	"github.com/Microsoft/go-winio"
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
	user32          = windows.NewLazySystemDLL("user32.dll")
	registerHotKey  = user32.NewProc("RegisterHotKey")
	getMessage      = user32.NewProc("GetMessageW")
	conn            net.Conn
	waitingForReply = false
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
	go startPipeServer()

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

func startPipeServer() {
	pipePath := `\\.\pipe\textwrench-pipe`
	listener, err := winio.ListenPipe(pipePath, nil)
	if err != nil {
		log.Fatalf("❌ Failed to create pipe: %v", err)
	}
	fmt.Println("📡 Named pipe server ready at", pipePath)

	for {
		c, err := listener.Accept()
		if err != nil {
			log.Println("❌ Accept error:", err)
			continue
		}
		conn = c
		fmt.Println("🔌 Electron connected")
		go handleClient(c)
	}
}

func handleClient(c net.Conn) {
	scanner := bufio.NewScanner(c)
	for scanner.Scan() {
		response := scanner.Text()
		fmt.Println("✅ Received processed:", response)
		onProcessedText(response)
	}
}

func handleHotkey() {
	if conn == nil {
		log.Println("⚠️ No connection to Electron")
		return
	}

	if waitingForReply {
		log.Println("⏳ Still waiting for previous reply")
		return
	}

	fmt.Println("🔁 Hotkey triggered")
	waitingForReply = true

	clipboard.WriteAll("") // Clear clipboard

	if err := sendCtrlKey(keybd_event.VK_C); err != nil {
		log.Println("❌ Failed to send Ctrl+C:", err)
		waitingForReply = false
		return
	}

	time.Sleep(100 * time.Millisecond)

	selected, err := clipboard.ReadAll()
	if err != nil || selected == "" {
		log.Println("❌ Could not read clipboard or empty selection")
		waitingForReply = false
		return
	}

	fmt.Println("📋 Selected:", selected)

	// Send selected text to Electron
	fmt.Fprintln(conn, selected)
}

func onProcessedText(text string) {
	defer func() { waitingForReply = false }()

	if err := clipboard.WriteAll(text); err != nil {
		log.Println("❌ Failed to write to clipboard:", err)
		return
	}

	time.Sleep(100 * time.Millisecond)

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
