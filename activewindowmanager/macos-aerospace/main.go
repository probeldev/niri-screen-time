// Package macosaerospace - implementation for MacOs AeroSpace
package macosaerospace

import (
	"encoding/json"
	"log"

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
		return appID, title, err
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

	log.Println(appID, title)

	return appID, title, nil
}
