package activewindowmanager

import (
	"errors"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/probeldev/niri-screen-time/activewindowmanager/hyprland"
	"github.com/probeldev/niri-screen-time/activewindowmanager/macos"
	macosaerospace "github.com/probeldev/niri-screen-time/activewindowmanager/macos-aerospace"
	"github.com/probeldev/niri-screen-time/activewindowmanager/niri"
	"github.com/probeldev/niri-screen-time/bash"
)

type ActiveWindowManagerInterface interface {
	GetActiveWindow() (string, string, error)
}

type CompositorType string

const (
	CompositorTypeNiri     CompositorType = "niri"
	CompositorTypeHyprland CompositorType = "hyprland"
)

func GetActiveWindowManager() (
	ActiveWindowManagerInterface,
	error,
) {
	currentOs := runtime.GOOS

	switch currentOs {
	case "darwin":
		return GetMacOsActiveWindowManager()
	case "linux":
		return GetLinuxActiveWindowManager()
	}

	return nil, errors.New("OS is not support")
}

func GetLinuxActiveWindowManager() (
	ActiveWindowManagerInterface,
	error,
) {
	desktop := os.Getenv("XDG_CURRENT_DESKTOP")
	if desktop == "" {
		return nil, errors.New("XDG_CURRENT_DESKTOP is not setup")
	}

	switch strings.ToLower(desktop) {
	case string(CompositorTypeNiri):
		manager := niri.NewNiriActiveWindow()
		return manager, nil
	case string(CompositorTypeHyprland):
		manager := hyprland.NewHyprlandActiveWindow()
		return manager, nil
	}

	return nil, errors.New("compositor is not supported")
}

func GetMacOsActiveWindowManager() (
	ActiveWindowManagerInterface,
	error,
) {

	if isSetCommand("aerospace -v") {
		log.Println("MacOs aerospace")
		return macosaerospace.NewMacOsAerospaceActiveWindow(), nil
	}
	log.Println("MacOs default")

	return macos.NewMacOsActiveWindow(), nil
}

func isSetCommand(command string) bool {
	_, err := bash.RunCommand(command)
	log.Println(err)
	return err == nil
}
