package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/probeldev/niri-screen-time/daemon"
	"github.com/probeldev/niri-screen-time/db"
	"github.com/probeldev/niri-screen-time/report"
)

func main() {
	db, err := db.NewScreenTimeDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	isDaemon := flag.Bool("daemon", false, "Run daemon")
	fromStr := flag.String("from", "", "Начальная дата (формат: 2006-01-02), по умолчанию — начало сегодняшнего дня")
	toStr := flag.String("to", "", "Конечная дата (формат: 2006-01-02), по умолчанию — конец сегодняшнего дня")

	flag.Parse()

	if *isDaemon {
		log.Println("Run daemon")
		daemon.Run(db)
	}

	flag.Parse()

	// Парсим даты (если не указаны — берём сегодняшний день)
	from, to := parseDates(*fromStr, *toStr)

	fmt.Println("")
	fmt.Printf("From %s to %s\n", from.Format("2006-01-02 15:04:05"), to.Format("2006-01-02 15:04:05"))

	err = report.GetReport(db, from, to)

	if err != nil {
		log.Panic(err)
	}
}

func parseDates(fromStr, toStr string) (from, to time.Time) {
	now := time.Now()

	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24*time.Hour - 1*time.Nanosecond)

	if fromStr == "" {
		from = todayStart
	} else {
		parsedFrom, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			log.Panic("Error parse date")
		} else {
			from = parsedFrom
		}
	}

	if toStr == "" {
		to = todayEnd
	} else {
		parsedTo, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			log.Panic("Error parse date")
		} else {
			to = parsedTo
		}
	}

	return from, to
}
