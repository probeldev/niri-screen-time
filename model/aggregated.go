package model

import (
	"time"
)

type AggregatedScreenTime struct {
	ID    int
	Date  time.Time
	AppID string
	Title string
	Sleep int
}

func NewAggregatedScreenTimeFromScreenTime(
	screenTime ScreenTime,
) AggregatedScreenTime {
	asc := AggregatedScreenTime{
		Date:  screenTime.Date,
		AppID: screenTime.AppID,
		Title: screenTime.Title,
		Sleep: screenTime.Sleep,
	}

	return asc
}

func (ast *AggregatedScreenTime) AddScreenTime(
	screenTime ScreenTime,
) {
	ast.Date = screenTime.Date
	ast.Sleep += screenTime.Sleep
}
