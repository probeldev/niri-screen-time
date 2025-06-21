package hyprland

import (
	"github.com/probeldev/niri-screen-time/bash"
	"strings"
)

type hyprlandActiveWindow struct{}

func NewHyprlandActiveWindow() *hyprlandActiveWindow {
	return &hyprlandActiveWindow{}
}

func (hyprlandActiveWindow) GetActiveWindow() (
	string,
	string,
	error,
) {

	appId := ""
	title := ""

	output, err := bash.RunCommand("hyprctl activewindow")
	if err != nil {
		return appId, title, err
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "class:") {
			appId = strings.TrimSpace(strings.TrimPrefix(line, "class:"))
		} else if strings.HasPrefix(line, "title:") {
			title = strings.TrimSpace(strings.TrimPrefix(line, "title:"))
		}
	}

	return appId, title, nil
}
