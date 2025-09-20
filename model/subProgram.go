package model

type SubProgram struct {
	AppIDs    []string `json:"app_ids"`
	TitleList []string `json:"title_list"`
	Alias     string   `json:"alias"`
}
