// Package daemon - a background service that collects data on application usage time (requires adding to startup for proper operation).
package daemon

import (
	"log"
	"time"

	"github.com/probeldev/niri-screen-time/activewindowmanager"
	"github.com/probeldev/niri-screen-time/cache"
	"github.com/probeldev/niri-screen-time/model"
)

const (
	sleepMs  = 200
	filename = "/home/sergey/screen-time.txt"
)

func Run(stc *cache.ScreenTimeCache) {
	fn := "daemon:Run"

	wm, err := activewindowmanager.GetActiveWindowManager()
	if err != nil {
		log.Panic(fn, err)
	}

	for {
		go func() {
			appID, title, err := wm.GetActiveWindow()
			if err != nil {
				log.Panic(fn, err)
			}

			if appID != "" {
				sc := model.ScreenTime{
					Date:  time.Now(),
					AppID: appID,
					Title: title,
					Sleep: sleepMs,
				}

				stc.Add(sc)
			}
		}()

		time.Sleep(sleepMs * time.Millisecond)
	}
}
