// Package responsemanager need for formatting response
package responsemanager

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
	"unicode/utf8"

	"github.com/probeldev/niri-screen-time/model"
)

type responseManagerCli struct {
	from  time.Time
	to    time.Time
	limit int
}

func NewResponseManagerCli(
	from time.Time,
	to time.Time,
	limit int,
) *responseManagerCli {
	r := responseManagerCli{}
	r.from = from
	r.to = to
	r.limit = limit

	return &r
}

func (r *responseManagerCli) Write(
	report []model.Report,
) {
	fn := "ResponseCli:write"

	fmt.Printf("\nReport period: %s to %s\n",
		r.from.Format("2006-01-02 15:04:05"),
		r.to.Format("2006-01-02 15:04:05"))
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

	var err error
	for i, rep := range report {
		if r.limit != 0 && i == r.limit {
			break
		}

		summary += rep.TimeMs
		dur := r.formatDuration(rep.TimeMs)

		name := strings.ReplaceAll(rep.Name, "\n", "")
		name = r.truncateString(name)
		_, err = fmt.Fprintf(w, "%s\t %s\n", name, dur)
		if err != nil {
			log.Println(fn, err)
		}
	}

	_, err = fmt.Fprintln(w, "\t\t")
	if err != nil {
		log.Println(fn, err)
	}

	dur := r.formatDuration(summary)
	_, err = fmt.Fprintf(w, "%s\t %s\n", "Summary screen time:", dur)
	if err != nil {
		log.Println(fn, err)
	}
	fmt.Println("")
}

func (*responseManagerCli) formatDuration(ms int) string {
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

func (*responseManagerCli) truncateString(s string) string {
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
