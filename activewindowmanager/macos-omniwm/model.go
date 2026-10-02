// Package macosomniwm - implementation for macOS OmniWM window manager
package macosomniwm

type App struct {
	BundleID string `json:"bundleId"`
	Name     string `json:"name"`
}

type Window struct {
	App       App    `json:"app"`
	Title     string `json:"title"`
	IsFocused bool   `json:"isFocused"`
}

type focusedWindowPayload struct {
	Window Window `json:"window"`
}

type focusedWindowResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		Kind    string               `json:"kind"`
		Payload focusedWindowPayload `json:"payload"`
	} `json:"result"`
}
