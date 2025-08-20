package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/probeldev/niri-screen-time/aggregatemanager"
	"github.com/probeldev/niri-screen-time/cache"
	"github.com/probeldev/niri-screen-time/daemon"
	"github.com/probeldev/niri-screen-time/db"
	"github.com/probeldev/niri-screen-time/details"
	"github.com/probeldev/niri-screen-time/report"
)

type Config struct {
	IsDaemon   bool
	IsDetails  bool
	From       string
	To         string
	AppID      string
	Title      string
	Limit      int
	IsOnlyText bool
}

func main() {
	if err := run(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := parseFlags()

	if cfg.IsDaemon {
		return runDaemonMode()
	}
	if cfg.IsDetails {
		return runDetailsMode(
			cfg.From,
			cfg.To,
			cfg.AppID,
			cfg.Title,
			cfg.Limit,
			cfg.IsOnlyText,
		)
	}
	return runReportMode(
		cfg.From,
		cfg.To,
	)
}

func parseFlags() *Config {
	cfg := &Config{}

	showVersion := false

	flag.BoolVar(&cfg.IsDaemon, "daemon", false, "Run daemon")
	flag.BoolVar(&cfg.IsDetails, "details", false, "View details")
	flag.BoolVar(&cfg.IsOnlyText, "onlytext", false, "Hack for remove counter from title")
	flag.StringVar(&cfg.From, "from", "", "Start date (format: 2006-01-02), defaults to today")
	flag.StringVar(&cfg.To, "to", "", "End date (format: 2006-01-02), defaults to today")
	flag.StringVar(&cfg.AppID, "appid", "", "AppId")
	flag.StringVar(&cfg.Title, "title", "", "Substring to match in titles")
	flag.IntVar(&cfg.Limit, "limit", 0, "Limit of response line, defaults to unlimited")
	flag.BoolVar(&showVersion, "version", false, "print version and exit")
	flag.Parse()

	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	return cfg
}

func runDaemonMode() error {
	fn := "runDaemonMode"
	// Создаем подключение к БД
	conn, err := db.NewDBConnection()
	if err != nil {
		log.Panic(fn, err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Panic(fn, err)
		}
	}()

	go func() {
		err := conn.Vacuum()
		if err != nil {
			log.Panic(fn, err)
		}
		time.Sleep(1 * time.Hour)
	}()

	if err := conn.InitTables(); err != nil {
		log.Panic(fn, err)
	}

	screenDB := db.NewScreenTimeDB(conn)
	aggregateDB := db.NewAggregatedScreenTimeDB(conn)

	go func() {
		am := aggregatemanager.NewAggragetManager(
			*screenDB,
			*aggregateDB,
		)

		am.Aggregate()
	}()

	screenTimeCache := cache.NewScreenTimeCache(screenDB, 5*time.Second, 100)
	screenTimeCache.Start()
	defer screenTimeCache.Stop()

	log.Println("Starting daemon...")

	daemon.Run(screenTimeCache)

	return nil
}

func runReportMode(
	fromStr string,
	toStr string,
) error {
	fn := "runReportMode"
	// Создаем подключение к БД
	conn, err := db.NewDBConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Panic(fn, err)
		}
	}()

	if err = conn.InitTables(); err != nil {
		log.Panic(fn, err)
	}

	screenTimeDB := db.NewScreenTimeDB(conn)

	from, to, err := parseDates(fromStr, toStr)
	if err != nil {
		return fmt.Errorf("failed to parse dates: %w", err)
	}

	fmt.Printf("\nReport period: %s to %s\n",
		from.Format("2006-01-02 15:04:05"),
		to.Format("2006-01-02 15:04:05"))

	aggregateDB := db.NewAggregatedScreenTimeDB(conn)

	return report.GetReport(
		screenTimeDB,
		aggregateDB,
		from,
		to,
	)
}

func runDetailsMode(
	fromStr string,
	toStr string,
	appID string,
	title string,
	limit int,
	isOnlyText bool,
) error {
	fn := "runDetailsMode"
	// Создаем подключение к БД
	conn, err := db.NewDBConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Panic(fn, err)
		}
	}()

	if err = conn.InitTables(); err != nil {
		log.Panic(fn, err)
	}

	screenTimeDB := db.NewScreenTimeDB(conn)

	from, to, err := parseDates(fromStr, toStr)
	if err != nil {
		return fmt.Errorf("failed to parse dates: %w", err)
	}

	fmt.Printf("\nReport period: %s to %s\n",
		from.Format("2006-01-02 15:04:05"),
		to.Format("2006-01-02 15:04:05"))

	aggregateDB := db.NewAggregatedScreenTimeDB(conn)

	return details.GetDetails(
		screenTimeDB,
		aggregateDB,
		from,
		to,
		appID,
		title,
		limit,
		isOnlyText,
	)
}

func parseDates(
	fromStr,
	toStr string,
) (
	from time.Time,
	to time.Time,
	err error,
) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24*time.Hour - 1*time.Nanosecond)

	parseDate := func(dateStr string, defaultDate time.Time) (time.Time, error) {
		if dateStr == "" {
			return defaultDate, nil
		}
		return time.Parse("2006-01-02", dateStr)
	}

	from, err = parseDate(fromStr, todayStart)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid from date: %w", err)
	}

	to, err = parseDate(toStr, todayEnd)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid to date: %w", err)
	}

	return from, to, nil
}
