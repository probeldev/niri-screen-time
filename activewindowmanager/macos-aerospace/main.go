// Package macosaerospace - implementation for MacOs AeroSpace
package macosaerospace

import (
	"encoding/json"

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
	output, err := bash.RunCommand("aerospace list-windows --focused --json")
	if err != nil {
		// We don't return the error because when there is no focused window,
		// the aerospace command exits with status code 1.
		return appID, title, nil
	}
	windows := []Window{}

	err = json.Unmarshal([]byte(output), &windows)
	if err != nil {
		return appID, title, err
	}

	if len(windows) == 0 {
		return "", "", nil
	}

	appID = windows[0].AppName
	title = windows[0].WindowTitle

	return appID, title, nil
}
