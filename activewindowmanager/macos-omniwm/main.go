// Package macosomniwm - implementation for macOS OmniWM window manager
package macosomniwm

import (
	"encoding/json"

	"github.com/probeldev/niri-screen-time/bash"
)

type macosOmniWMActiveWindow struct{}

func NewMacOsOmniWMActiveWindow() *macosOmniWMActiveWindow {
	return &macosOmniWMActiveWindow{}
}

func (macosOmniWMActiveWindow) GetActiveWindow() (
	appID string,
	title string,
	err error,
) {
	output, err := bash.RunCommand("omniwmctl query focused-window")
	if err != nil {
		// OmniWM may not be running, its IPC server may be disabled
		// (general.ipcEnabled), or there may be no focused window.
		return appID, title, nil
	}

	response := focusedWindowResponse{}

	err = json.Unmarshal([]byte(output), &response)
	if err != nil {
		return appID, title, err
	}

	window, ok := parseFocusedWindow(response)
	if !ok {
		return appID, title, nil
	}

	appID = window.App.Name
	if appID == "" {
		appID = window.App.BundleID
	}
	title = window.Title

	return appID, title, nil
}

// parseFocusedWindow извлекает окно из ответа OmniWM. Пустой payload —
// штатная ситуация: сфокусированного управляемого окна нет (например,
// фокус на рабочем столе).
func parseFocusedWindow(response focusedWindowResponse) (Window, bool) {
	if !response.OK {
		return Window{}, false
	}

	window := response.Result.Payload.Window
	if window.App.Name == "" && window.App.BundleID == "" {
		return Window{}, false
	}

	return window, true
}
