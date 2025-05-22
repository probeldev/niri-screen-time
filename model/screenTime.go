package model

import "time"

type ScreenTime struct {
	ID    int
	Date  time.Time
	AppID string
	Title string
	Sleep int
}
