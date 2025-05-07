package hotkey

import (
	"context"
	"fmt"

	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

type HotkeyManager struct {
	hk *hotkey.Hotkey
}

// NewHotkeyManager returns a hotkey manager ready to register hotkeys.
func NewHotkeyManager() *HotkeyManager {
	return &HotkeyManager{}
}

// Register sets up a hotkey and listens for key events in the background.
func (h *HotkeyManager) Register(ctx context.Context, mods []hotkey.Modifier, key hotkey.Key, callback func()) error {
	h.hk = hotkey.New(mods, key)
	if err := h.hk.Register(); err != nil {
		return fmt.Errorf("failed to register hotkey: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-h.hk.Keydown():
				callback()
			}
		}
	}()
	return nil
}

func (h *HotkeyManager) UnregisterAll() {
	if h.hk != nil {
		h.hk.Unregister()
	}
}

// InitMainThread wraps mainthread.Init, which is required for some OSes.
func InitMainThread(mainFunc func()) {
	mainthread.Init(mainFunc)
}
