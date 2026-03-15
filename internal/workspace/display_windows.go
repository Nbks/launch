//go:build windows

package workspace

import (
	"fmt"
	"launch/internal/logger"
	"syscall"
	"unsafe"
)

type Display struct {
	Name   string
	Width  int
	Height int
	X      int
	Y      int
}

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procEnumDisplayMonitors = user32.NewProc("EnumDisplayMonitors")
	procGetMonitorInfo      = user32.NewProc("GetMonitorInfoW")
)

const (
	MONITORINFOF_PRIMARY = 0x00000001
)

type RECT struct {
	Left, Top, Right, Bottom int32
}

type MONITORINFOEX struct {
	CbSize    uint32
	RcMonitor RECT
	RcWork    RECT
	DwFlags   uint32
	SzDevice  [32]uint16 // CCHDEVICENAME = 32
}

func GetDisplays() ([]Display, error) {
	var displays []Display
	var monitorIndex int = 0

	callback := syscall.NewCallback(func(hMonitor uintptr, hdc uintptr, lprcMonitor uintptr, dwData uintptr) uintptr {
		var mi MONITORINFOEX
		mi.CbSize = uint32(unsafe.Sizeof(mi))

		// Get monitor info
		ret, _, _ := procGetMonitorInfo.Call(
			hMonitor,
			uintptr(unsafe.Pointer(&mi)),
		)

		if ret == 0 {
			logger.Log.Error("Failed to get monitor info")
			return 1 // Continue enumeration
		}

		// Calculate dimensions
		w := int(mi.RcMonitor.Right - mi.RcMonitor.Left)
		h := int(mi.RcMonitor.Bottom - mi.RcMonitor.Top)
		x := int(mi.RcMonitor.Left)
		y := int(mi.RcMonitor.Top)

		// Determine display name
		// If primary monitor, name it DISPLAY0, otherwise DISPLAY<n> where n is incrementing
		displayName := fmt.Sprintf("DISPLAY%d", monitorIndex)

		if mi.DwFlags&MONITORINFOF_PRIMARY != 0 {
			displayName = "DISPLAY0"
		}

		monitorIndex++

		display := Display{
			Name:   displayName,
			Width:  w,
			Height: h,
			X:      x,
			Y:      y,
		}

		displays = append(displays, display)
		logger.Log.Debug("Monitor detected", "display", display)
		return 1 // Continue enumeration
	})

	// Enumerate all monitors
	ret, _, _ := procEnumDisplayMonitors.Call(0, 0, callback, 0)
	if ret == 0 {
		return nil, fmt.Errorf("failed to enumerate monitors")
	}

	// If no displays were found, return an error
	if len(displays) == 0 {
		return nil, fmt.Errorf("no displays detected")
	}

	return displays, nil
}
