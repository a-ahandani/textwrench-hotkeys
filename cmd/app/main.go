package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"textwrench-hotkeys/internal/clipboard"
	"textwrench-hotkeys/internal/comms"
	"textwrench-hotkeys/internal/hotkey"
	"time"
)

type ShortcutConfig struct {
	ID        string   `json:"id"`
	Key       string   `json:"key"`
	Modifiers []string `json:"modifiers"`
}

func main() {
	hotkey.InitMainThread(run)
}

func run() {
	// Create cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Instantiate communicator and hotkey manager
	comm := comms.NewCommunicator()
	manager := hotkey.NewManager()

	// Graceful shutdown handling: on SIGINT/SIGTERM, close comm and cancel context
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("\nShutting down...")
		if err := comm.Close(); err != nil {
			fmt.Printf("Failed to close communicator: %v\n", err)
		}
		cancel()
	}()

	// Register default hotkeys
	defaultHotkeys := []ShortcutConfig{
		{ID: "fixSelectedText", Key: "C", Modifiers: []string{"ctrl", "shift"}},
		{ID: "explainSelectedText", Key: "E", Modifiers: []string{"ctrl", "shift"}},
		{ID: "selectPrompt", Key: "J", Modifiers: []string{"ctrl", "shift"}},
	}
	manager.UnregisterAll()
	registerHotkeys(ctx, manager, comm, defaultHotkeys)

	// Message handler from the front-end
	handler := func(message string) {
		if strings.HasPrefix(message, "SHORTCUT_CONFIG|") {
			handleConfigMessage(message, comm, manager, ctx)
			return
		}

		fmt.Printf("Received message at %s: %s\n", time.Now().Format(time.RFC3339), message)
		switch {
		case strings.HasPrefix(message, "FOCUS_PASTE|"):
			content := strings.TrimPrefix(message, "FOCUS_PASTE|")
			if err := clipboard.FocusPasteText(content); err != nil {
				fmt.Printf("Failed to focus and paste text: %v\n", err)
			}

		case strings.HasPrefix(message, "PASTE|"):
			content := strings.TrimPrefix(message, "PASTE|")

			if err := clipboard.WriteText(content); err != nil {
				fmt.Printf("Failed to paste text: %v\n", err)
			}

		default:
			// Backward‐compatible fallback
			fmt.Printf("Ignoring unrecognized comm message: %q\n", message)
		}
	}

	// Start the IPC listener
	go func() {
		if err := comm.Start(ctx, handler); err != nil {
			log.Fatalf("Failed to start communicator: %v", err)
		}
	}()

	fmt.Println("Hotkey app running. Waiting for configuration...")
	<-ctx.Done()

	// Final cleanup (in case signal handler didn't already)
	if err := comm.Close(); err != nil {
		fmt.Printf("Failed to close communicator: %v\n", err)
	}
	manager.UnregisterAll()
	fmt.Println("Exited neatly")
}

func registerHotkeys(ctx context.Context, manager *hotkey.Manager, comm comms.Communicator, hotkeys []ShortcutConfig) {
	for _, cfg := range hotkeys {
		id := cfg.ID
		key := parseKey(cfg.Key)
		mods := parseModifiers(cfg.Modifiers)

		action := func() {
			text, err := clipboard.ReadSelectedText()
			if err != nil || strings.TrimSpace(text) == "" {
				return
			}
			comm.Send(fmt.Sprintf("%s|%s", id, text))
		}

		if err := manager.RegisterHotkey(ctx, id, mods, key, action); err != nil {
			fmt.Printf("Failed to register hotkey %s: %v\n", id, err)
		}
	}
}

func handleConfigMessage(message string, comm comms.Communicator, manager *hotkey.Manager, ctx context.Context) {
	configJSON := strings.TrimPrefix(message, "SHORTCUT_CONFIG|")
	var configs []ShortcutConfig
	if err := json.Unmarshal([]byte(configJSON), &configs); err != nil {
		fmt.Printf("Invalid config received: %v\n", err)
		return
	}
	fmt.Println("Registering new shortcuts configuration")
	manager.UnregisterAll()
	registerHotkeys(ctx, manager, comm, configs)
}

func parseModifiers(mods []string) []hotkey.Modifier {
	var result []hotkey.Modifier
	for _, mod := range mods {
		switch strings.ToLower(mod) {
		case "ctrl":
			result = append(result, hotkey.ModCtrl)
		case "shift":
			result = append(result, hotkey.ModShift)
		case "alt", "option":
			result = append(result, hotkey.ModShift) // fallback, adjust if you define `ModAlt`
		case "cmd", "meta":
			result = append(result, hotkey.ModCmd)
		}
	}
	return result
}

func parseKey(k string) hotkey.Key {
	switch strings.ToUpper(k) {
	case "A":
		return hotkey.KeyA
	case "B":
		return hotkey.KeyB
	case "C":
		return hotkey.KeyC
	case "D":
		return hotkey.KeyD
	case "E":
		return hotkey.KeyE
	case "F":
		return hotkey.KeyF
	case "G":
		return hotkey.KeyG
	case "H":
		return hotkey.KeyH
	case "I":
		return hotkey.KeyI
	case "J":
		return hotkey.KeyJ
	case "K":
		return hotkey.KeyK
	case "L":
		return hotkey.KeyL
	case "M":
		return hotkey.KeyM
	case "N":
		return hotkey.KeyN
	case "O":
		return hotkey.KeyO
	case "P":
		return hotkey.KeyP
	case "Q":
		return hotkey.KeyQ
	case "R":
		return hotkey.KeyR
	case "S":
		return hotkey.KeyS
	case "T":
		return hotkey.KeyT
	case "U":
		return hotkey.KeyU
	case "V":
		return hotkey.KeyV
	case "W":
		return hotkey.KeyW
	case "X":
		return hotkey.KeyX
	case "Y":
		return hotkey.KeyY
	case "Z":
		return hotkey.KeyZ
	default:
		return hotkey.KeyC
	}
}
