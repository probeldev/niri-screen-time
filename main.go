package main

import (
	"flag"
	"log"
	"time"

	"github.com/probeldev/niri-screen-time/daemon"
	"github.com/probeldev/niri-screen-time/db"
	"github.com/probeldev/niri-screen-time/report"
)

func main() {
	db, err := db.NewScreenTimeDB("screen_time.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	isDaemon := flag.Bool("daemon", false, "Run daemon")

	flag.Parse()

	if *isDaemon {
		log.Println("Run daemon")
		daemon.Run(db)
	}

	from := time.Now().Add(-1 * 24 * 365 * time.Hour)
	to := time.Now()

	err = report.GetReport(db, from, to)

	if err != nil {
		log.Panic(err)
	}
}
