// Package macos - implementation for macos
package macos

import (
	"strings"

	"github.com/probeldev/niri-screen-time/bash"
)

type macosActiveWindow struct{}

func NewMacOsActiveWindow() *macosActiveWindow {
	return &macosActiveWindow{}
}

func (macosActiveWindow) GetActiveWindow() (
	appID string,
	title string,
	err error,
) {
	output, err := bash.RunCommand(`
		osascript <<EOF
		tell application "System Events"
			set frontApp to first application process whose frontmost is true
			set appName to name of frontApp
			set windowName to name of window 1 of frontApp
			set processID to unix id of frontApp
			return "App: " & appName & " | Window: " & windowName & " | PID: " & processID
		end tell
		EOF
	`)
	if err != nil {
		return appID, title, err
	}

	lines := strings.Split(output, "|")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if after, found := strings.CutPrefix(line, "App:"); found {
			appID = strings.TrimSpace(after)
		} else if after, found := strings.CutPrefix(line, "Window:"); found {
			title = strings.TrimSpace(after)
		}
	}

	return appID, title, nil
}
