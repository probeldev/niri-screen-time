package macosrift

type WindowID struct {
	PID int `json:"pid"`
	Idx int `json:"idx"`
}

type Window struct {
	ID             WindowID `json:"id"`
	Title          string   `json:"title"`
	IsFloating     bool     `json:"is_floating"`
	IsFocused      bool     `json:"is_focused"`
	BundleID       string   `json:"bundle_id"`
	AppName        string   `json:"app_name"`
	WindowServerID int      `json:"window_server_id"`
}
