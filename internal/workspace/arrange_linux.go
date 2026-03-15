package workspace

import (
	"fmt"
	"os/exec"
	"time"

	"launch/internal/config"
	"strings"
)

func ApplyLayout(t config.Tool) error {
	if t.Layout == "" {
		return nil
	}
	fmt.Println("ENTRO EN APPLAYOUT", t)

	displays, err := GetDisplays()
	if err != nil {
		return err
	}

	fmt.Println(displays)
	var target *Display

	for _, d := range displays {
		if d.Name == t.Display {
			target = &d
			break
		}
	}

	if target == nil {
		return fmt.Errorf("display '%s' not found", t.Display)
	}

	fmt.Println(t)
	// Esperar a que la ventana exista
	err = waitForWindow(t.Name, 10*time.Second)
	if err != nil {
		fmt.Println("Window not detected:", err)
		return err
	}
	switch t.Layout {

	case "fullscreen":
		return exec.Command("wmctrl",
			"-r", t.Name,
			"-b", "add,maximized_vert,maximized_horz",
		).Run()
	case "left-half":
		return moveWindow(t.Name, target.X, target.Y, target.Width/2, target.Height)

	case "right-half":
		return moveWindow(t.Name, target.X+target.Width/2, target.Y, target.Width/2, target.Height)

	case "top-half":
		return moveWindow(t.Name, target.X, target.Y, target.Width, target.Height/2)

	case "bottom-half":
		return moveWindow(t.Name, target.X, target.Y+target.Height/2, target.Width, target.Height/2)

	default:
		return nil
	}
}

func moveWindow(name string, x, y, w, h int) error {
	return exec.Command("wmctrl",
		"-r", name,
		"-e", fmt.Sprintf("0,%d,%d,%d,%d", x, y, w, h),
	).Run()
}

func waitForWindow(name string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {

		out, _ := exec.Command("wmctrl", "-l").Output()

		if strings.Contains(string(out), name) {
			return nil
		}

		time.Sleep(200 * time.Millisecond)
	}

	return fmt.Errorf("window '%s' not found after %v", name, timeout)
}
