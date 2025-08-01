// Package hyprland - implementation for hyprland wayland conpositor
package hyprland

import (
	"strings"

	"github.com/probeldev/niri-screen-time/bash"
)

type hyprlandActiveWindow struct{}

func NewHyprlandActiveWindow() *hyprlandActiveWindow {
	return &hyprlandActiveWindow{}
}

func (hyprlandActiveWindow) GetActiveWindow() (
	appID string,
	title string,
	err error,
) {
	output, err := bash.RunCommand("hyprctl activewindow")
	if err != nil {
		return appID, title, err
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if after, found := strings.CutPrefix(line, "class:"); found {
			appID = strings.TrimSpace(after)
		} else if after, found := strings.CutPrefix(line, "title:"); found {
			title = strings.TrimSpace(after)
		}
	}

	return appID, title, nil
}
