// Package macosrift - implementation for macOS Rift window manager
package macosrift

import (
	"encoding/json"

	"github.com/probeldev/niri-screen-time/bash"
)

type macosRiftActiveWindow struct{}

func NewMacOsRiftActiveWindow() *macosRiftActiveWindow {
	return &macosRiftActiveWindow{}
}

func (macosRiftActiveWindow) GetActiveWindow() (
	appID string,
	title string,
	err error,
) {
	output, err := bash.RunCommand("rift-cli query windows")
	if err != nil {
		// Rift may not be running or the query may fail when there is no
		// focused window.
		return appID, title, nil
	}

	windows := []Window{}

	err = json.Unmarshal([]byte(output), &windows)
	if err != nil {
		return appID, title, err
	}

	for _, window := range windows {
		if window.IsFocused {
			appID = window.AppName
			if appID == "" {
				appID = window.BundleID
			}
			title = window.Title
			return appID, title, nil
		}
	}

	return "", "", nil
}
