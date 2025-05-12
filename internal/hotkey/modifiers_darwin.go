//go:build darwin
// +build darwin

package hotkey

import "golang.design/x/hotkey"

const (
	ModCtrl  = hotkey.ModCtrl
	ModShift = hotkey.ModShift
	ModCmd   = hotkey.ModCmd
	ModMeta  = hotkey.ModOption
)
