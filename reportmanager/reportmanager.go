// Package reportmanager need for write report by using apps
package reportmanager

import (
	"log"
	"time"

	"github.com/probeldev/niri-screen-time/aliasmanager"
	"github.com/probeldev/niri-screen-time/db"
	"github.com/probeldev/niri-screen-time/model"
	"github.com/probeldev/niri-screen-time/subprogrammanager"
)

type ResponseManagerInterface interface {
	Write([]model.Report)
}

type reportManager struct {
	responseManager ResponseManagerInterface
}

func NewResponseManager(
	responseManager ResponseManagerInterface,
) *reportManager {
	r := reportManager{}
	r.responseManager = responseManager

	return &r
}

func (r *reportManager) GetReport(
	dbScreenTime *db.ScreenTimeDB,
	dbAggregate *db.AggregatedScreenTimeDB,
	from time.Time,
	to time.Time,
) error {
	resp := map[string]model.Report{}

	screenTimeList, err := dbScreenTime.GetByDateRange(
		from,
		to,
	)

	if err != nil {
		return err
	}

	aggregate, err := dbAggregate.GetByDateRange(
		from,
		to,
	)
	if err != nil {
		return err
	}

	for _, a := range aggregate {
		screenTimeList = append(screenTimeList, model.ScreenTime(a))
	}

	summary := 0
	subProgram, err := subprogrammanager.NewSubProgramManager()

	if err != nil {
		return err
	}

	for _, st := range screenTimeList {
		summary += st.Sleep
		st = subProgram.GetSubProgram(st)

		if report, ok := resp[st.AppID]; ok {
			report.TimeMs += st.Sleep
			resp[st.AppID] = report
		} else {
			resp[st.AppID] = model.Report{
				Name:   st.AppID,
				TimeMs: st.Sleep,
			}
		}
	}

	responseSlice := []model.Report{}
	for _, responseApp := range resp {
		responseSlice = append(responseSlice, responseApp)
	}

	responseSlice, err = r.useAliace(responseSlice)
	if err != nil {
		return err
	}

	r.responseManager.Write(responseSlice)

	return nil
}

func (r *reportManager) useAliace(reports []model.Report) ([]model.Report, error) {
	fn := "reportManager:useAliace"

	alias, err := aliasmanager.NewAliasManager()
	if err != nil {
		log.Println(fn, err)
		return nil, err
	}

	reportsWithAlias := []model.Report{}
	for _, r := range reports {
		alias := alias.ReplaceAppId2Alias(r)
		reportsWithAlias = append(reportsWithAlias, alias)
	}

	return reportsWithAlias, nil
}
