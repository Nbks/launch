package workspace

import (
	"fmt"
	"os/exec"
	"strings"
)

type Display struct {
	Name   string
	Width  int
	Height int
	X      int
	Y      int
}

func GetDisplays() ([]Display, error) {
	fmt.Println("---- [workspace] Running xrandr --query ----")
	out, err := exec.Command("xrandr", "--query").Output()
	if err != nil {
		fmt.Println("error", err)
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	var displays []Display

	for _, line := range lines {
		if strings.Contains(line, " connected") && !strings.Contains(line, "disconnected") {
			var name string
			var resolution string

			parts := strings.Fields(line)
			name = parts[0]

			for _, p := range parts {
				if strings.Contains(p, "+") && strings.Contains(p, "x") {
					resolution = p
					break
				}
			}

			if resolution == "" {
				continue
			}

			var w, h, x, y int
			fmt.Sscanf(resolution, "%dx%d+%d+%d", &w, &h, &x, &y)

			displays = append(displays, Display{
				Name:   name,
				Width:  w,
				Height: h,
				X:      x,
				Y:      y,
			})
		}
	}

	return displays, nil
}
