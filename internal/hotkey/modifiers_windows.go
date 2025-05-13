//go:build windows

package hotkey

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"golang.design/x/hotkey"
	hk "golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

type Modifier = hotkey.Modifier
type Key = hotkey.Key

const (
	KeyA = hotkey.KeyA
	KeyB = hotkey.KeyB
	KeyC = hotkey.KeyC
	KeyD = hotkey.KeyD
	KeyE = hotkey.KeyE
	KeyF = hotkey.KeyF
	KeyG = hotkey.KeyG
	KeyH = hotkey.KeyH
	KeyI = hotkey.KeyI
	KeyJ = hotkey.KeyJ
	KeyK = hotkey.KeyK
	KeyL = hotkey.KeyL
	KeyM = hotkey.KeyM
	KeyN = hotkey.KeyN
	KeyO = hotkey.KeyO
	KeyP = hotkey.KeyP
	KeyQ = hotkey.KeyQ
	KeyR = hotkey.KeyR
	KeyS = hotkey.KeyS
	KeyT = hotkey.KeyT
	KeyU = hotkey.KeyU
	KeyV = hotkey.KeyV
	KeyW = hotkey.KeyW
	KeyX = hotkey.KeyX
	KeyY = hotkey.KeyY
	KeyZ = hotkey.KeyZ
)

type registeredHotkey struct {
	instance *hotkey.Hotkey
	cancel   context.CancelFunc
}

type Manager struct {
	mu       sync.Mutex
	registry map[string]*registeredHotkey
}

// InitMainThread wraps required mainthread initialization
func InitMainThread(fn func()) {
	mainthread.Init(fn)
}

// NewManager creates a hotkey manager
func NewManager() *Manager {
	return &Manager{
		registry: make(map[string]*registeredHotkey),
	}
}

func (m *Manager) RegisterHotkey(ctx context.Context, id string, mods []Modifier, key Key, cb func()) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Remove existing if present
	if existing, ok := m.registry[id]; ok {
		existing.instance.Unregister()
		existing.cancel()
		delete(m.registry, id)
	}

	hk := hotkey.New(mods, key)
	if err := hk.Register(); err != nil {
		return fmt.Errorf("failed to register hotkey: %w", err)
	}

	hkCtx, cancel := context.WithCancel(ctx)
	go func() {
		for {
			select {
			case <-hkCtx.Done():
				return
			case <-hk.Keydown():
				cb()
			}
		}
	}()

	m.registry[id] = &registeredHotkey{instance: hk, cancel: cancel}
	return nil
}

// UnregisterAll removes all active hotkeys
func (m *Manager) UnregisterAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for id, hk := range m.registry {
		hk.instance.Unregister()
		hk.cancel()
		delete(m.registry, id)
	}
}

// ParseModifiers maps strings like "ctrl", "shift", "cmd" to Windows-specific modifiers
func ParseModifiers(mods []string) []Modifier {
	var result []Modifier
	for _, mod := range mods {
		switch strings.ToLower(mod) {
		case "ctrl":
			result = append(result, Modifier(hk.ModCtrl))
		case "shift":
			result = append(result, Modifier(hk.ModShift))
		case "cmd", "meta":
			result = append(result, Modifier(hk.ModCtrl)) // Treat as Ctrl on Windows
		}
	}
	return result
}
