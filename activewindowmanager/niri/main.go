package niri

import (
	niriwindows "github.com/probeldev/niri-float-sticky/niri-windows"
)

type niriActiveWindow struct{}

func NewNiriActiveWindow() *niriActiveWindow {
	return &niriActiveWindow{}
}

func (niriActiveWindow) GetActiveWindow() (
	string,
	string,
	error,
) {
	windows, err := niriwindows.GetWindowsList()

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
