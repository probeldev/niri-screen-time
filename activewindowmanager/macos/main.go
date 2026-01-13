// Package macos - implementation for macos
package macos

import (
	"fmt"
	"strings"

	"github.com/probeldev/niri-screen-time/bash"
)

func RequestPermissions() error {
	script := `
	osascript <<'EOF'
	tell application "System Events"
		if not UI elements enabled then
			display dialog "This app requires Accessibility permissions to track active windows." & return & return & ¬
			"Please click 'Open Settings' and add this app to the allowed list." ¬
			with title "Permission Required" ¬
			with icon caution ¬
			buttons {"Open Settings", "Cancel"} ¬
			default button 1
			if button returned of result is "Open Settings" then
				do shell script "open 'x-apple.systempreferences:com.apple.preference.security?Privacy_Accessibility'"
			end if
		else
			display dialog "Accessibility permissions are already enabled!" ¬
			buttons {"OK"} ¬
			default button 1
		end if
	end tell
	EOF
	`

	_, err := bash.RunCommand(script)
	return err
}

type macosActiveWindow struct {
	permissionsChecked bool
}

func NewMacOsActiveWindow() *macosActiveWindow {
	return &macosActiveWindow{}
}

func (m *macosActiveWindow) CheckPermissions() error {
	script := `
	osascript -e 'tell application "System Events" to get UI elements enabled'
	`

	output, err := bash.RunCommand(script)
	if err != nil {
		return fmt.Errorf("failed to check permissions: %v", err)
	}

	if strings.TrimSpace(output) == "false" {
		return fmt.Errorf("accessibility permissions not granted")
	}

	m.permissionsChecked = true
	return nil
}

func (m *macosActiveWindow) GetActiveWindow() (appID string, title string, err error) {
	// Проверяем permissions перед выполнением
	if !m.permissionsChecked {
		if errPemission := m.CheckPermissions(); errPemission != nil {
			return "", "", fmt.Errorf("permissions required: %v", errPemission)
		}
	}

	output, err := bash.RunCommand(`
		osascript <<'EOF'
		try
			tell application "System Events"
				set frontApp to first application process whose frontmost is true
				set appName to name of frontApp
				set processID to unix id of frontApp
				
				try
					set windowName to name of window 1 of frontApp
				on error
					set windowName to "No Window"
				end try
				
				return "App: " & appName & " | Window: " & windowName & " | PID: " & processID
			end tell
		on error errMsg
			return "Error: " & errMsg
		end try
		EOF
	`)
	if err != nil {
		return "", "", fmt.Errorf("failed to get active window: %v", err)
	}

	// Обрабатываем ошибки из AppleScript
	if strings.Contains(output, "Error:") {
		return "", "", fmt.Errorf("apple script error: %s", strings.TrimPrefix(output, "Error: "))
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

// EnsurePermissions проверяет и запрашивает права если нужно
func (m *macosActiveWindow) EnsurePermissions() error {
	if err := m.CheckPermissions(); err != nil {
		fmt.Println("🔧 Требуются права доступа...")
		return RequestPermissions()
	}
	return nil
}
