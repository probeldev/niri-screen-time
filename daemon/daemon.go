package daemon

import (
	"log"
	"time"

	niriwindows "github.com/probeldev/niri-float-sticky/niri-windows"
	"github.com/probeldev/niri-screen-time/cache"
	"github.com/probeldev/niri-screen-time/model"
)

const (
	sleepMs  = 200
	filename = "/home/sergey/screen-time.txt"
)

func Run(cache *cache.ScreenTimeCache) {

	fn := "daemon:Run"

	for {

		go func() {
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

					cache.Add(sc)
				}
			}
		}()

		time.Sleep(sleepMs * time.Millisecond)
	}
}
