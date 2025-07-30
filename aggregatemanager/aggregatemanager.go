package aggregatemanager

import (
	"log"
	"time"

	"github.com/probeldev/niri-screen-time/db"
	"github.com/probeldev/niri-screen-time/model"
)

type aggregateManager struct {
	screenTimeDb db.ScreenTimeDB
	aggregateDb  db.AggregatedScreenTimeDB
}

func NewAggragetManager(
	screenTimeDb db.ScreenTimeDB,
	aggregateDb db.AggregatedScreenTimeDB,
) *aggregateManager {
	am := &aggregateManager{}
	am.screenTimeDb = screenTimeDb
	am.aggregateDb = aggregateDb

	return am
}

func (am *aggregateManager) Aggregate() {
	for {
		am.aggregateWorker()
		time.Sleep(10 * time.Minute)
	}
}

func (am *aggregateManager) aggregateWorker() {
	fn := "aggregateManager:aggregateWorker"
	screenTimes, err := am.screenTimeDb.GetAll()

	if err != nil {
		log.Println(fn, err)
		return
	}

	var aggregate model.AggregatedScreenTime
	screenTimeForDelete := []model.ScreenTime{}

	for _, st := range screenTimes {
		if len(screenTimeForDelete) == 0 {
			aggregate = model.NewAggregatedScreenTimeFromScreenTime(st)

			screenTimeForDelete = append(screenTimeForDelete, st)
			continue
		}

		if am.needAggregate(aggregate, st) {
			aggregate.AddScreenTime(st)
			screenTimeForDelete = append(screenTimeForDelete, st)
			continue
		}

		err := am.aggregateDb.Insert(aggregate)
		if err != nil {
			log.Println(fn, err)
			return
		}

		for _, std := range screenTimeForDelete {
			err := am.screenTimeDb.DeleteByID(std)
			if err != nil {
				log.Println(fn, err)
				return
			}
		}

		aggregate = model.NewAggregatedScreenTimeFromScreenTime(st)
		screenTimeForDelete = []model.ScreenTime{st}
	}

	log.Println("finished")
}

func (*aggregateManager) needAggregate(
	aggregate model.AggregatedScreenTime,
	screenTime model.ScreenTime,
) bool {
	if aggregate.AppID != screenTime.AppID {
		return false
	}

	if aggregate.Title != screenTime.Title {
		return false
	}

	if aggregate.Date.Sub(screenTime.Date) > time.Second {
		return false
	}

	if aggregate.Date.Format("2006-01-02") != screenTime.Date.Format("2006-01-02") {
		return false
	}

	return true
}
