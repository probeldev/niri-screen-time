// Package details write details for using apps
package details

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
	"unicode/utf8"

	"github.com/probeldev/niri-screen-time/aliasmanager"
	"github.com/probeldev/niri-screen-time/db"
	"github.com/probeldev/niri-screen-time/model"
)

func GetDetails(
	dbScreenTime *db.ScreenTimeDB,
	dbAggregate *db.AggregatedScreenTimeDB,
	from time.Time,
	to time.Time,
	appID string,
	title string,
	limit int,
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

	for _, st := range screenTimeList {
		summary += st.Sleep

		if st.AppID != appID {
			continue
		}

		if !strings.Contains(st.Title, title) {
			continue
		}

		if report, ok := response[st.Title]; ok {
			report.TimeMs += st.Sleep
			response[st.Title] = report
		} else {
			response[st.Title] = model.Report{
				Name:   st.Title,
				TimeMs: st.Sleep,
			}
		}
	}

	responseSlice := []model.Report{}
	for _, responseApp := range response {
		responseSlice = append(responseSlice, responseApp)
	}

	write(responseSlice, limit)

	return nil
}

// TODO: move to new package
func write(report []model.Report, limit int) {
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

	// TODO: тут не нужны алиасы
	for i, r := range report {
		if limit != 0 && i == limit {
			break
		}
		summary += r.TimeMs
		dur := formatDuration(r.TimeMs)

		alias := alias.ReplaceAppId2Alias(r)

		name := truncateString(alias.Name)
		_, err = fmt.Fprintf(w, "%s\t %s\n", name, dur)
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

// TODO: move to new package
func formatDuration(ms int) string {
	if ms < 0 {
		return "0ms"
	}

	seconds := ms / 1000
	minutes := seconds / 60
	hours := minutes / 60

	seconds %= 60
	minutes %= 60

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

func truncateString(s string) string {
	maxLength := 80
	// Если строка короче или равна максимальной длине, возвращаем как есть
	if utf8.RuneCountInString(s) <= maxLength {
		return s
	}

	// Преобразуем строку в срез рун для корректной работы с Unicode
	runes := []rune(s)

	// Обрезаем до maxLength символов и добавляем многоточие
	return string(runes[:maxLength]) + "..."
}
