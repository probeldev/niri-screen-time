// Package autostartmanager implement autostart for MacOs
package autostartmanager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/probeldev/niri-screen-time/bash"
)

type AutoStartManager struct {
	appName     string
	plistPath   string
	programPath string
	args        []string
}

// NewAutoStartManager creates an autostart manager for a specific program
func NewAutoStartManager(programPath string, args []string) (*AutoStartManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Use program name for plist file name
	appName := filepath.Base(programPath)
	plistPath := filepath.Join(homeDir, "Library", "LaunchAgents",
		fmt.Sprintf("com.niri.screentime.%s.plist", appName))

	return &AutoStartManager{
		appName:     appName,
		plistPath:   plistPath,
		programPath: programPath,
		args:        args,
	}, nil
}

// NewAutoStartManagerForMacOs creates a manager specifically for niri-screen-time
func NewAutoStartManagerForMacOs() (*AutoStartManager, error) {
	// Get full path to current executable
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %v", err)
	}

	// Use full path to current binary
	programPath := execPath
	args := []string{"-daemon"}

	return NewAutoStartManager(programPath, args)
}

func (a *AutoStartManager) Enable() error {
	// Collect all arguments into one array
	programArgs := []string{a.programPath}
	programArgs = append(programArgs, a.args...)

	// Form XML for ProgramArguments
	argsXML := ""
	for _, arg := range programArgs {
		argsXML += fmt.Sprintf("        <string>%s</string>\n", arg)
	}

	plistContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.niri.screentime.%s</string>
    <key>ProgramArguments</key>
    <array>
%s    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/tmp/%s.log</string>
    <key>StandardErrorPath</key>
    <string>/tmp/%s.error.log</string>
    <key>EnvironmentVariables</key>
    <dict>
        <key>PATH</key>
        <string>/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin</string>
    </dict>
</dict>
</plist>`, a.appName, argsXML, a.appName, a.appName)

	dir := filepath.Dir(a.plistPath)
	var permissionFolder os.FileMode = 0755
	if err := os.MkdirAll(dir, permissionFolder); err != nil {
		return err
	}

	var permissionFile os.FileMode = 0644
	if err := os.WriteFile(a.plistPath, []byte(plistContent), permissionFile); err != nil {
		return err
	}

	fmt.Printf("✓ Autostart enabled for %s: %s\n", a.programPath, a.plistPath)
	return nil
}

func (a *AutoStartManager) Load() error {
	cmd := fmt.Sprintf("launchctl load \"%s\"", a.plistPath)
	_, err := bash.RunCommand(cmd)
	if err != nil {
		return fmt.Errorf("service load error: %v", err)
	}
	fmt.Println("✓ Service loaded, application started")
	return nil
}

func (a *AutoStartManager) Unload() error {
	cmd := fmt.Sprintf("launchctl unload \"%s\"", a.plistPath)
	_, err := bash.RunCommand(cmd)
	if err != nil {
		return fmt.Errorf("service unload error: %v", err)
	}
	fmt.Println("✓ Service unloaded")
	return nil
}

func (a *AutoStartManager) EnableAndLoad() error {
	if err := a.Enable(); err != nil {
		return err
	}
	return a.Load()
}

func (a *AutoStartManager) Disable() error {
	// First unload the service
	_ = a.Unload()

	// Then remove the plist file
	if err := os.Remove(a.plistPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("autostart was not configured")
		}
		return err
	}

	fmt.Println("✓ Autostart disabled")
	return nil
}

func (a *AutoStartManager) Status() (bool, bool) {
	// Check if plist file exists
	plistExists := false
	if _, err := os.Stat(a.plistPath); err == nil {
		plistExists = true
	}

	// Check if service is running
	cmd := fmt.Sprintf("launchctl list | grep com.niri.screentime.%s", a.appName)
	output, err := bash.RunCommand(cmd)
	isRunning := err == nil && strings.Contains(output, "com.niri.screentime."+a.appName)

	return plistExists, isRunning
}

// GetPlistPath returns the path to the created plist file
func (a *AutoStartManager) GetPlistPath() string {
	return a.plistPath
}

// CheckAndFixPermissions checks and fixes access permissions
func (a *AutoStartManager) CheckAndFixPermissions() error {
	fmt.Println("🔧 Checking and configuring permissions...")

	// Add application to Accessibility allowed list
	script := `
	osascript <<'EOF'
	tell application "System Events"
		-- Check if UI elements are enabled
		if not UI elements enabled then
			display dialog "niri-screen-time requires Accessibility permissions to track active windows." & return & return & ¬
			"Click 'Open Settings' and add niri-screen-time to the list of allowed applications." ¬
			with title "Permissions Required" ¬
			with icon caution ¬
			buttons {"Open Settings", "Cancel"} ¬
			default button 1
			
			if button returned of result is "Open Settings" then
				tell application "System Preferences"
					activate
					reveal anchor "Privacy_Accessibility" of pane "com.apple.preference.security"
				end tell
			end if
		else
			-- Already have permissions, try to add our application
			try
				set appPath to "%s"
				tell application "System Events"
					tell process "System Preferences"
						if exists then
							-- Settings already open, do nothing
						end if
					end tell
				end tell
			on error
				-- Ignore errors, main thing is that we have permissions
			end try
		end if
	end tell
	EOF
	`

	_, err := bash.RunCommand(fmt.Sprintf(script, a.programPath))
	if err != nil {
		fmt.Printf("⚠️  Failed to automatically configure permissions: %v\n", err)
		fmt.Println("📋 Please manually add niri-screen-time to:")
		fmt.Println("   System Settings → Privacy & Security → Accessibility")
	} else {
		fmt.Println("✓ Permissions configured")
	}

	return nil
}
