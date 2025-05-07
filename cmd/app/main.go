package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"textwrench-hotkeys/internal/clipboard"
	"textwrench-hotkeys/internal/comms"
	"textwrench-hotkeys/internal/hotkey"

	hk "golang.design/x/hotkey"
)

func main() {
	hotkey.InitMainThread(run)
}

func run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("\nShutting down...")
		cancel()
	}()

	// Initialize communication module (pipe/socket)
	comm := comms.NewCommunicator()
	handler := func(message string) {
		fmt.Printf("Received message: %s\n", message)

		// Paste the message (simulate Ctrl+V or Cmd+V)
		err := clipboard.PasteText(message)
		if err != nil {
			fmt.Printf("Failed to paste text: %v\n", err)
		}
	}

	go func() {
		if err := comm.Start(ctx, handler); err != nil {
			log.Fatalf("Failed to start communicator: %v", err)
		}
	}()

	// Initialize hotkey manager
	manager := hotkey.NewHotkeyManager()
	fmt.Println("App started")

	err := manager.Register(ctx, []hk.Modifier{hk.ModCtrl, hk.ModShift}, hk.KeyC, func() {
		text, err := clipboard.ReadSelectedText()
		fmt.Println("Hotkey triggered: Ctrl+Shift+C", text)
		if err != nil {
			fmt.Printf("Error reading clipboard: %v\n", err)
			return
		}

		err = comm.Send(text)
		if err != nil {
			fmt.Printf("Error sending text via comms: %v\n", err)
		}
	})

	if err != nil {
		fmt.Printf("Failed to register hotkey: %v\n", err)
		return
	}

	fmt.Println("Hotkey registered. Waiting for events...")
	<-ctx.Done()
	manager.UnregisterAll()
	fmt.Println("Exited cleanly")
}
