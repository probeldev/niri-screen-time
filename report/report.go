package report

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/probeldev/niri-screen-time/aliasmanager"
	"github.com/probeldev/niri-screen-time/db"
	"github.com/probeldev/niri-screen-time/model"
	"github.com/probeldev/niri-screen-time/subprogrammanager"
)

func GetReport(
	dbScreenTime *db.ScreenTimeDB,
	dbAggregate *db.AggregatedScreenTimeDB,
	from time.Time,
	to time.Time,
) error {
	response := map[string]model.Report{}

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
	subProgram, err := subprogrammanager.NewSubProgramManager()

	if err != nil {
		return err
	}

	for _, st := range screenTimeList {
		summary += st.Sleep
		st = subProgram.GetSubProgram(st)

		if report, ok := response[st.AppID]; ok {
			report.TimeMs += st.Sleep
			response[st.AppID] = report
		} else {
			response[st.AppID] = model.Report{
				Name:   st.AppID,
				TimeMs: st.Sleep,
			}
		}
	}

	responseSlice := []model.Report{}
	for _, responseApp := range response {
		responseSlice = append(responseSlice, responseApp)
	}

	write(responseSlice)

	return nil
}

func write(report []model.Report) {
	fn := "report:write"

	sort.Slice(report, func(i, j int) bool {
		return report[i].TimeMs > report[j].TimeMs
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer func() {
		err := w.Flush()
		if err != nil {
			log.Panic(fn, err)
		}
	}()

	summary := 0

	alias, err := aliasmanager.NewAliasManager()
	if err != nil {
		log.Panic(fn, err)
	}

	for _, r := range report {
		summary += r.TimeMs
		dur := formatDuration(r.TimeMs)

		alias := alias.ReplaceAppId2Alias(r)
		_, err = fmt.Fprintf(w, "%s\t %s\n", alias.Name, dur)
		if err != nil {
			log.Println(fn, err)
		}
	}

	_, err = fmt.Fprintln(w, "\t\t")
	if err != nil {
		log.Println(fn, err)
	}

	dur := formatDuration(summary)
	_, err = fmt.Fprintf(w, "%s\t %s\n", "Summary screen time:", dur)
	if err != nil {
		log.Println(fn, err)
	}
	fmt.Println("")
}

func formatDuration(ms int) string {
	if ms < 0 {
		return "0ms"
	}

	seconds := ms / 1000
	minutes := seconds / 60
	hours := minutes / 60

	seconds = seconds % 60
	minutes = minutes % 60

	parts := []string{}

	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	if seconds > 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}
	if ms%1000 > 0 && hours == 0 {
		parts = append(parts, fmt.Sprintf("%dms", ms%1000))
	}

	if len(parts) == 0 {
		return "0ms"
	}

	return strings.Join(parts, " ")
}
