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

func Run(cache *cache.ScreenTimeCache) {
	fn := "daemon:Run"

	wm, err := activewindowmanager.GetActiveWindowManager()
	if err != nil {
		log.Panic(fn, err)
	}

	for {
		go func() {
			appId, title, err := wm.GetActiveWindow()
			if err != nil {
				log.Panic(fn, err)
			}

			if appId != "" {
				sc := model.ScreenTime{
					Date:  time.Now(),
					AppID: appId,
					Title: title,
					Sleep: sleepMs,
				}

				cache.Add(sc)
			}
		}()

		time.Sleep(sleepMs * time.Millisecond)
	}
}
