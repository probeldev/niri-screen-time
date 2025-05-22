package model

import "time"

type AggregatedScreenTime struct {
	Date  time.Time
	AppID string
	Title string
	Sleep int
}
