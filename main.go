package main

import (
	"log"

	"github.com/probeldev/niri-screen-time/db"
	"github.com/probeldev/niri-screen-time/demon"
)

func main() {
	db, err := db.NewScreenTimeDB("screen_time.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	demon.Run(db)
}
