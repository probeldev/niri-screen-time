// Package niri. Realization for wayland compositor Niri
package niri

import (
	"encoding/json"
	"fmt"

	"github.com/probeldev/niri-screen-time/bash"
)

type niriActiveWindow struct{}

func NewNiriActiveWindow() *niriActiveWindow {
	return &niriActiveWindow{}
}

func (nw *niriActiveWindow) GetWindowsList() ([]Window, error) {
	output, err := bash.RunCommand("niri msg --json windows")
	if err != nil {
		return nil, err
	}

	windows, err := nw.ParseWindows([]byte(output))

	return windows, err
}

func (*niriActiveWindow) ParseWindows(output []byte) ([]Window, error) {
	var windows []Window
	if err := json.Unmarshal(output, &windows); err != nil {
		return nil, fmt.Errorf("error unmarshalling windows: %w", err)
	}
	return windows, nil
}

func (nw *niriActiveWindow) GetActiveWindow() (
	appID string,
	title string,
	err error,
) {
	windows, err := nw.GetWindowsList()

	if err != nil {
		return "", "", err
	}

	for _, w := range windows {
		if w.IsFocused {
			return w.AppID, w.Title, nil
		}
	}

	return "", "", nil
}
