//go:build windows

package workspace

import (
	"fmt"
	"launch/internal/config"
	"launch/internal/logger"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var (
	// Import user32 functions needed for window manipulation
	procEnumWindows     = user32.NewProc("EnumWindows")
	procGetWindowText   = user32.NewProc("GetWindowTextW")
	procIsWindowVisible = user32.NewProc("IsWindowVisible")
	procSetWindowPos    = user32.NewProc("SetWindowPos")
	procShowWindow      = user32.NewProc("ShowWindow")
	procFindWindow      = user32.NewProc("FindWindowW")
	procGetWindowRect   = user32.NewProc("GetWindowRect")
)

const (
	SW_MAXIMIZE       = 3
	SW_RESTORE        = 9
	SWP_NOACTIVATE    = 0x0010
	SWP_NOOWNERZORDER = 0x0200
	SWP_SHOWWINDOW    = 0x0040
	HWND_TOP          = 0
)

// Window handle storage for the waitForWindow function
type windowFind struct {
	hwnd    uintptr
	title   string
	pattern string
	found   bool
}

func ApplyLayout(t config.Tool) error {
	logger.Log.Debug("Applying layout", "tool", t.Name, "layout", t.Layout)

	if t.Layout == "" {
		return nil
	}

	// Get display information
	displays, err := GetDisplays()
	if err != nil {
		logger.Log.Error("Failed to get displays", "error", err)
		return err
	}

	// Find the target display
	var target *Display
	for _, d := range displays {
		if d.Name == t.Display {
			target = &d
			break
		}
	}

	// If display not specified, use the first one (which should be primary)
	if t.Display == "" && len(displays) > 0 {
		target = &displays[0]
	} else if target == nil {
		return fmt.Errorf("display '%s' not found", t.Display)
	}

	// Wait for window to appear
	hwnd, err := waitForWindow(t.Name, 10*time.Second)
	if err != nil {
		logger.Log.Error("Failed to find window", "window", t.Name, "error", err)
		return err
	}

	logger.Log.Debug("Found window", "window", t.Name, "hwnd", hwnd)

	// Apply the requested layout
	switch t.Layout {
	case "fullscreen":
		return setFullscreen(hwnd)
	case "left-half":
		return moveWindow(hwnd, target.X, target.Y, target.Width/2, target.Height)
	case "right-half":
		return moveWindow(hwnd, target.X+target.Width/2, target.Y, target.Width/2, target.Height)
	case "top-half":
		return moveWindow(hwnd, target.X, target.Y, target.Width, target.Height/2)
	case "bottom-half":
		return moveWindow(hwnd, target.X, target.Y+target.Height/2, target.Width, target.Height/2)
	default:
		logger.Log.Warn("Unknown layout type", "layout", t.Layout)
		return nil
	}
}

// setFullscreen maximizes the window
func setFullscreen(hwnd uintptr) error {
	ret, _, err := procShowWindow.Call(hwnd, uintptr(SW_MAXIMIZE))
	if ret == 0 {
		return fmt.Errorf("failed to maximize window: %v", err)
	}
	return nil
}

// moveWindow repositions and resizes a window
func moveWindow(hwnd uintptr, x, y, width, height int) error {
	// First restore the window if it's maximized
	procShowWindow.Call(hwnd, uintptr(SW_RESTORE))

	// Set window position and size
	ret, _, err := procSetWindowPos.Call(
		hwnd,
		uintptr(HWND_TOP),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(SWP_NOACTIVATE|SWP_NOOWNERZORDER|SWP_SHOWWINDOW),
	)

	if ret == 0 {
		return fmt.Errorf("failed to set window position: %v", err)
	}
	return nil
}

// waitForWindow waits for a window with the given title pattern to appear
func waitForWindow(pattern string, timeout time.Duration) (uintptr, error) {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		// Try to find the window
		finder := &windowFind{
			pattern: pattern,
			found:   false,
		}

		// Enumerate all windows
		procEnumWindows.Call(
			syscall.NewCallback(func(hwnd uintptr, lparam uintptr) uintptr {
				// Check if window is visible
				isVisible, _, _ := procIsWindowVisible.Call(hwnd)
				if isVisible == 0 {
					return 1 // Continue enumeration
				}

				// Get window title
				var title [256]uint16
				length, _, _ := procGetWindowText.Call(
					hwnd,
					uintptr(unsafe.Pointer(&title[0])),
					uintptr(len(title)),
				)

				if length > 0 {
					windowTitle := syscall.UTF16ToString(title[:length])
					logger.Log.Debug("Found window", "title", windowTitle)

					// Check if the window title contains our pattern
					if strings.Contains(strings.ToLower(windowTitle), strings.ToLower(pattern)) {
						finder.hwnd = hwnd
						finder.title = windowTitle
						finder.found = true
						return 0 // Stop enumeration
					}
				}
				return 1 // Continue enumeration
			}),
			0,
		)

		if finder.found {
			return finder.hwnd, nil
		}

		time.Sleep(200 * time.Millisecond)
	}

	return 0, fmt.Errorf("window with title containing '%s' not found after %v", pattern, timeout)
}

// ListWindows is a helper function for debugging that returns a list of all visible windows
// This can be used to help identify window titles for testing purposes
func ListWindows() []string {
	var windows []string

	procEnumWindows.Call(
		syscall.NewCallback(func(hwnd uintptr, lparam uintptr) uintptr {
			// Check if window is visible
			isVisible, _, _ := procIsWindowVisible.Call(hwnd)
			if isVisible == 0 {
				return 1 // Continue enumeration
			}

			// Get window title
			var title [256]uint16
			length, _, _ := procGetWindowText.Call(
				hwnd,
				uintptr(unsafe.Pointer(&title[0])),
				uintptr(len(title)),
			)

			if length > 0 {
				windowTitle := syscall.UTF16ToString(title[:length])
				if windowTitle != "" {
					windows = append(windows, windowTitle)
				}
			}
			return 1 // Continue enumeration
		}),
		0,
	)

	return windows
}
