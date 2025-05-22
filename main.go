package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/probeldev/niri-screen-time/cache"
	"github.com/probeldev/niri-screen-time/daemon"
	"github.com/probeldev/niri-screen-time/db"
	"github.com/probeldev/niri-screen-time/report"
)

type Config struct {
	IsDaemon bool
	From     string
	To       string
}

func main() {
	if err := run(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := parseFlags()
	if err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	if cfg.IsDaemon {
		return runDaemonMode()
	}
	return runReportMode(cfg.From, cfg.To)
}

func parseFlags() (*Config, error) {
	cfg := &Config{}

	flag.BoolVar(&cfg.IsDaemon, "daemon", false, "Run daemon")
	flag.StringVar(&cfg.From, "from", "", "Start date (format: 2006-01-02), defaults to today")
	flag.StringVar(&cfg.To, "to", "", "End date (format: 2006-01-02), defaults to today")
	flag.Parse()

	return cfg, nil
}

func runDaemonMode() error {

	// Создаем подключение к БД
	conn, err := db.NewDBConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err := conn.InitTables(); err != nil {
		log.Fatal(err)
	}

	db := db.NewScreenTimeDB(conn)

	cache := cache.NewScreenTimeCache(db, 5*time.Second, 100)
	cache.Start()
	defer cache.Stop()

	log.Println("Starting daemon...")

	daemon.Run(cache)

	return nil
}

func runReportMode(fromStr, toStr string) error {
	// Создаем подключение к БД
	conn, err := db.NewDBConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err := conn.InitTables(); err != nil {
		log.Fatal(err)
	}

	db := db.NewScreenTimeDB(conn)

	from, to, err := parseDates(fromStr, toStr)
	if err != nil {
		return fmt.Errorf("failed to parse dates: %w", err)
	}

	fmt.Printf("\nReport period: %s to %s\n",
		from.Format("2006-01-02 15:04:05"),
		to.Format("2006-01-02 15:04:05"))

	return report.GetReport(db, from, to)
}

func parseDates(fromStr, toStr string) (time.Time, time.Time, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24*time.Hour - 1*time.Nanosecond)

	parseDate := func(dateStr string, defaultDate time.Time) (time.Time, error) {
		if dateStr == "" {
			return defaultDate, nil
		}
		return time.Parse("2006-01-02", dateStr)
	}

	from, err := parseDate(fromStr, todayStart)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid from date: %w", err)
	}

	to, err := parseDate(toStr, todayEnd)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid to date: %w", err)
	}

	return from, to, nil
}
