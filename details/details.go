// Package details write details for using apps
package details

import (
	"strconv"
	"strings"
	"time"

	"github.com/probeldev/niri-screen-time/db"
	"github.com/probeldev/niri-screen-time/model"
	"github.com/probeldev/niri-screen-time/response"
)

func GetDetails(
	dbScreenTime *db.ScreenTimeDB,
	dbAggregate *db.AggregatedScreenTimeDB,
	from time.Time,
	to time.Time,
	appID string,
	title string,
	limit int,
	isOnlyText bool,
) error {
	resp := map[string]model.Report{}

	screenTimeList, err := dbScreenTime.GetByDateRange(
		from,
		to,
	)

	if err != nil {
		return err
	}

	aggregate, err := dbAggregate.GetByDateRange(
		from,
		to,
	)
	if err != nil {
		return err
	}

	for _, a := range aggregate {
		screenTimeList = append(screenTimeList, model.ScreenTime(a))
	}

	summary := 0

	for _, st := range screenTimeList {
		summary += st.Sleep

		if st.AppID != appID {
			continue
		}

		if !strings.Contains(st.Title, title) {
			continue
		}

		if isOnlyText {
			st.Title = onlyText(st.Title)
		}

		if report, ok := resp[st.Title]; ok {
			report.TimeMs += st.Sleep
			resp[st.Title] = report
		} else {
			resp[st.Title] = model.Report{
				Name:   st.Title,
				TimeMs: st.Sleep,
			}
		}
	}

	responseSlice := []model.Report{}
	for _, responseApp := range resp {
		responseSlice = append(responseSlice, responseApp)
	}

	response.Write(responseSlice, limit)

	return nil
}

func onlyText(s string) string {
	for i := range 10 {
		s = strings.ReplaceAll(s, strconv.Itoa(i), "")
	}

	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")
	s = strings.ReplaceAll(s, "–", "")
	s = strings.ReplaceAll(s, "-", "")
	s = strings.Trim(s, " ")

	return s
}
