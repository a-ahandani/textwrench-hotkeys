//go:build windows
// +build windows

package hotkey

import "golang.design/x/hotkey"

const (
	ModCtrl  = hotkey.ModCtrl
	ModShift = hotkey.ModShift
	ModCmd   = hotkey.ModCtrl
	ModAlt   = hotkey.ModAlt
)
