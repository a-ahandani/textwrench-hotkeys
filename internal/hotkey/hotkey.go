package hotkey

import "context"

// Modifier and Key types are declared in platform-specific files

// Manager defines the interface shared across platforms
type ManagerInterface interface {
	RegisterHotkey(ctx context.Context, id string, mods []Modifier, key Key, cb func()) error
	UnregisterAll()
}
