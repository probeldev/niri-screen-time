// Package macosomniwm - implementation for macOS OmniWM window manager
package macosomniwm

import (
	"encoding/json"
	"testing"
)

// Снимок реального ответа `omniwmctl query focused-window` (OmniWM 0.7.3).
const focusedWindowJSON = `{
  "id": "7F5AB558-93F2-44B6-9182-5659E48A0CE1",
  "kind": "query",
  "ok": true,
  "result": {
    "kind": "focused-window",
    "payload": {
      "window": {
        "app": {
          "bundleId": "net.kovidgoyal.kitty",
          "name": "kitty"
        },
        "isFocused": true,
        "mode": "tiling",
        "pid": 2603,
        "title": "OC | What is new in rift wm",
        "windowId": 130
      }
    }
  },
  "status": "success",
  "version": 16
}`

// OmniWM отвечает пустым payload, когда сфокусированного управляемого окна нет.
const noFocusedWindowJSON = `{
  "id": "7E01F93B-6AAE-4C58-B985-1151B0B7C8B3",
  "kind": "query",
  "ok": true,
  "result": {
    "kind": "focused-window",
    "payload": {}
  },
  "status": "success",
  "version": 16
}`

func TestParseFocusedWindow(t *testing.T) {
	response := focusedWindowResponse{}

	if err := json.Unmarshal([]byte(focusedWindowJSON), &response); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	window, ok := parseFocusedWindow(response)
	if !ok {
		t.Fatal("expected a focused window")
	}

	if window.App.Name != "kitty" {
		t.Errorf("app name = %q, want %q", window.App.Name, "kitty")
	}

	if window.Title != "OC | What is new in rift wm" {
		t.Errorf("title = %q, want %q", window.Title, "OC | What is new in rift wm")
	}
}

func TestParseFocusedWindowEmptyPayload(t *testing.T) {
	response := focusedWindowResponse{}

	if err := json.Unmarshal([]byte(noFocusedWindowJSON), &response); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if _, ok := parseFocusedWindow(response); ok {
		t.Fatal("expected no focused window for an empty payload")
	}
}
