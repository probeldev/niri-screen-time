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
	"github.com/probeldev/niri-screen-time/detailsmanager"
	"github.com/probeldev/niri-screen-time/model"
	"github.com/probeldev/niri-screen-time/reportmanager"
	"github.com/probeldev/niri-screen-time/responsemanager"
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
	IsJSON     bool
}

func main() {
	if err := run(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}

type ResponseManagerInterface interface {
	Write([]model.Report)
}

func GetResponseManager(
	cfg *Config,
) ResponseManagerInterface {

	fn := "GetResponseManager"

	if cfg.IsJSON {
		responseManager := responsemanager.NewResponseManagerJSON(
			cfg.Limit,
		)
		return responseManager
	}

	// TODO: move to parsing flags
	from, to, err := parseDates(cfg.From, cfg.To)
	if err != nil {
		log.Println(fn, err)
		os.Exit(0)
	}

	responseManager := responsemanager.NewResponseManagerCli(
		from,
		to,
		cfg.Limit,
	)
	return responseManager
}

func run() error {
	cfg := parseFlags()

	if cfg.IsDaemon {
		return runDaemonMode()
	}

	responseManager := GetResponseManager(cfg)

	if cfg.IsDetails {
		return runDetailsMode(
			cfg.From,
			cfg.To,
			cfg.AppID,
			cfg.Title,
			cfg.IsOnlyText,
			responseManager,
		)
	}
	return runReportMode(
		cfg.From,
		cfg.To,
		responseManager,
	)
}

func parseFlags() *Config {
	cfg := &Config{}

	showVersion := false

	flag.BoolVar(&cfg.IsDaemon, "daemon", false, "Run daemon")
	flag.BoolVar(&cfg.IsDetails, "details", false, "View details")
	flag.BoolVar(&cfg.IsOnlyText, "onlytext", false, "Hack for remove counter from title")
	flag.BoolVar(&cfg.IsJSON, "json", false, "return response with json format")
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
	responseManager reportmanager.ResponseManagerInterface,
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

	// TODO: move to parsing flags
	from, to, err := parseDates(fromStr, toStr)
	if err != nil {
		return fmt.Errorf("failed to parse dates: %w", err)
	}

	aggregateDB := db.NewAggregatedScreenTimeDB(conn)

	report := reportmanager.NewResponseManager(
		responseManager,
	)

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
	isOnlyText bool,
	responseManager detailsmanager.ResponseManagerInterface,
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

	// TODO: move to parsing flags
	from, to, err := parseDates(fromStr, toStr)
	if err != nil {
		return fmt.Errorf("failed to parse dates: %w", err)
	}

	aggregateDB := db.NewAggregatedScreenTimeDB(conn)

	details := detailsmanager.NewDealsManager(
		responseManager,
	)

	return details.GetDetails(
		screenTimeDB,
		aggregateDB,
		from,
		to,
		appID,
		title,
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
