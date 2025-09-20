// Package responsemanager need for formatting response
package responsemanager

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/probeldev/niri-screen-time/model"
)

type responseManagerJSON struct {
	limit int
}

func NewResponseManagerJSON(
	limit int,
) *responseManagerJSON {
	r := responseManagerJSON{}
	r.limit = limit

	return &r
}

func (r *responseManagerJSON) Write(
	report []model.Report,
) {
	fn := "ResponseJSON:write"

	sort.Slice(report, func(i, j int) bool {
		return report[i].TimeMs > report[j].TimeMs
	})

	jsonReport, err := json.Marshal(report)

	if err != nil {
		log.Println(fn, err)
		os.Exit(0)
	}

	fmt.Println(string(jsonReport))
}
