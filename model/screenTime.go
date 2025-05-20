package model

import "time"

type ScreenTime struct {
	Date  time.Time
	AppID string
	Title string
	Sleep int
}
