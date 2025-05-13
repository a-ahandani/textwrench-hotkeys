//go:build darwin

package hotkey

import (
	"context"
	"errors"
	"strings"
)

type Modifier int
type Key int

const (
	KeyA = iota
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
)

type Manager struct{}

// InitMainThread is a no-op on macOS for now
func InitMainThread(fn func()) {
	fn()
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) RegisterHotkey(ctx context.Context, id string, mods []Modifier, key Key, cb func()) error {
	return errors.New("hotkey registration not implemented on macOS")
}

func (m *Manager) UnregisterAll() {}

func ParseModifiers(mods []string) []Modifier {
	var result []Modifier
	for _, mod := range mods {
		switch strings.ToLower(mod) {
		case "ctrl":
			result = append(result, Modifier(1))
		case "shift":
			result = append(result, Modifier(2))
		case "cmd", "meta":
			result = append(result, Modifier(3))
		}
	}
	return result
}
