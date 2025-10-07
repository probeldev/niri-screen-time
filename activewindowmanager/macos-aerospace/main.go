// Package macosaerospace - implementation for MacOs AeroSpace
package macosaerospace

import (
	"strings"

	"github.com/probeldev/niri-screen-time/bash"
)

type macosAeropaceActiveWindow struct{}

func NewMacOsAerospaceActiveWindow() *macosAeropaceActiveWindow {
	return &macosAeropaceActiveWindow{}
}

func (macosAeropaceActiveWindow) GetActiveWindow() (
	appID string,
	title string,
	err error,
) {
	output, err := bash.RunCommand("aerospace list-windows --focused")
	if err != nil {
		return appID, title, err
	}

	if strings.Contains(output, "No window is focused") {
		return "", "", nil
	}

	lines := strings.Split(output, "|")
	appID = strings.Trim(lines[1], " ")
	title = strings.Trim(lines[2], " ")

	return appID, title, nil
}
