package main

import (
	"log"

	"github.com/probeldev/niri-screen-time/daemon"
	"github.com/probeldev/niri-screen-time/db"
)

func main() {
	db, err := db.NewScreenTimeDB("screen_time.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	daemon.Run(db)
}
