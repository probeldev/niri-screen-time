package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/probeldev/niri-screen-time/activewindowmanager/macos"
	"github.com/probeldev/niri-screen-time/aggregatemanager"
	"github.com/probeldev/niri-screen-time/autostartmanager"
	"github.com/probeldev/niri-screen-time/cache"
	"github.com/probeldev/niri-screen-time/daemon"
	"github.com/probeldev/niri-screen-time/db"
	"github.com/probeldev/niri-screen-time/detailsmanager"
	"github.com/probeldev/niri-screen-time/model"
	"github.com/probeldev/niri-screen-time/reportmanager"
	"github.com/probeldev/niri-screen-time/responsemanager"
)

type Config struct {
	IsDaemon       bool
	IsDetails      bool
	From           string
	To             string
	AppID          string
	Title          string
	Limit          int
	IsOnlyText     bool
	IsJSON         bool
	IsMacOsStartup bool
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

	if cfg.IsMacOsStartup {
		return manageAutoStart(cfg)
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
	flag.BoolVar(&cfg.IsMacOsStartup, "autostart", false, "manage macos autostart (enable/disable/status)")
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

func addToStartupMacOs() error {
	fmt.Println("🚀 Настройка автозапуска для macOS...")

	// Создаем менеджер автозапуска
	manager, err := autostartmanager.NewAutoStartManagerForMacOs()
	if err != nil {
		return fmt.Errorf("ошибка создания менеджера автозапуска: %v", err)
	}

	// Проверяем и настраиваем права
	windowTracker := macos.NewMacOsActiveWindow()
	if err := windowTracker.EnsurePermissions(); err != nil {
		fmt.Printf("⚠️  Предупреждение: %v\n", err)
		fmt.Println("📋 Для полной функциональности потребуются права доступа")
	}

	// Настраиваем права для launchd
	if err := manager.CheckAndFixPermissions(); err != nil {
		fmt.Printf("⚠️  Предупреждение: %v\n", err)
	}

	// Включаем автозапуск
	if err := manager.EnableAndLoad(); err != nil {
		return fmt.Errorf("ошибка включения автозапуска: %v", err)
	}

	// Проверяем статус
	plistExists, isRunning := manager.Status()
	fmt.Printf("\n📊 Статус автозапуска:\n")
	fmt.Printf("   Plist файл: %s\n", manager.GetPlistPath())
	fmt.Printf("   Plist существует: %t\n", plistExists)
	fmt.Printf("   Служба запущена: %t\n", isRunning)

	if isRunning {
		fmt.Println("\n✅ niri-screen-time успешно добавлен в автозапуск и запущен!")
	} else {
		fmt.Println("\n⚠️  Автозапуск настроен, но служба не запущена")
		fmt.Println("   Приложение запустится при следующей перезагрузке")
	}

	fmt.Println("\n💡 Для отключения автозапуска используйте: niri-screen-time autostart disable")

	return nil
}

func manageAutoStart(cfg *Config) error {
	// TODO: only darwin

	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("  niri-screen-time -autostart enable   - add to autostart")
		fmt.Println("  niri-screen-time -autostart disable  - remove from autostart")
		fmt.Println("  niri-screen-time -autostart status   - check status")
		return nil
	}

	manager, err := autostartmanager.NewAutoStartManagerForMacOs()
	if err != nil {
		return err
	}

	switch os.Args[2] {
	case "enable":
		return addToStartupMacOs()
	case "disable":
		return manager.Disable()
	case "status":
		plistExists, isRunning := manager.Status()
		fmt.Printf("Plist exists: %t\n", plistExists)
		fmt.Printf("Service is running: %t\n", isRunning)
		if plistExists {
			fmt.Printf("Plist path: %s\n", manager.GetPlistPath())
		}
	default:
		return fmt.Errorf("unknown command: %s", os.Args[2])
	}

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

	details := detailsmanager.NewDetailsManager(
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
