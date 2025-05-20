package daemon

import (
	"log"
	"time"

	niriwindows "github.com/probeldev/niri-float-sticky/niri-windows"
	"github.com/probeldev/niri-screen-time/db"
	"github.com/probeldev/niri-screen-time/model"
)

const (
	sleepMs  = 200
	filename = "/home/sergey/screen-time.txt"
)

func Run(db *db.ScreenTimeDB) {

	fn := "daemon:Run"

	for {
		windows, err := niriwindows.GetWindowsList()

		if err != nil {
			log.Panic(fn, err)
		}

		for _, w := range windows {
			if w.IsFocused {

				sc := model.ScreenTime{
					Date:  time.Now(),
					AppID: w.AppID,
					Title: w.Title,
					Sleep: sleepMs,
				}

				if err := db.Insert(sc); err != nil {
					log.Fatal(err)
				}

			}
		}

		time.Sleep(sleepMs * time.Millisecond)
	}
}
