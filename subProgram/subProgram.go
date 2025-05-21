package subprogram

import (
	"strings"

	"github.com/probeldev/niri-screen-time/model"
)

type subProgram struct {
	subProgram []model.SubProgram
}

func NewSubProgram() subProgram {
	sp := subProgram{}

	sp.subProgram = append(sp.subProgram, model.SubProgram{
		AppID: "com.mitchellh.ghostty",
		Title: "nvim",
		Alias: "NeoVim",
	})

	sp.subProgram = append(sp.subProgram, model.SubProgram{
		AppID: "zen",
		Title: "DeepSeek",
		Alias: "DeepSeek",
	})

	sp.subProgram = append(sp.subProgram, model.SubProgram{
		AppID: "zen",
		Title: "YouTube",
		Alias: "YouTube",
	})

	sp.subProgram = append(sp.subProgram, model.SubProgram{
		AppID: "org.telegram.desktop",
		Title: "ProBelDev Chat",
		Alias: "ProBelDev Chat",
	})
	return sp
}

func (sp *subProgram) isSetProgram(st model.ScreenTime) bool {
	for _, sp := range sp.subProgram {
		if sp.AppID == st.AppID {
			return true
		}
	}

	return false
}

func (sp *subProgram) GetSubProgram(st model.ScreenTime) model.ScreenTime {
	if !sp.isSetProgram(st) {
		return st
	}

	for _, sp := range sp.subProgram {
		if sp.AppID == st.AppID {
			if strings.Contains(st.Title, sp.Title) {
				st.AppID = st.AppID + "(" + sp.Alias + ")"
				return st
			}
		}
	}

	st.AppID = st.AppID + "(Other)"
	return st
}
