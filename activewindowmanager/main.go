package activewindowmanager

import (
	"errors"
	"github.com/probeldev/niri-screen-time/activewindowmanager/hyprland"
	"github.com/probeldev/niri-screen-time/activewindowmanager/niri"
	"os"
	"strings"
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
