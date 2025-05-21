package alias

import (
	"strings"

	"github.com/probeldev/niri-screen-time/model"
)

type alias struct {
	aliases []model.Alias
}

func NewAlias() alias {
	a := alias{}

	a.aliases = append(a.aliases, model.Alias{Name: "com.mitchellh.ghostty", Alias: "Ghostty"})
	a.aliases = append(a.aliases, model.Alias{Name: "org.telegram.desktop", Alias: "Telegram"})
	a.aliases = append(a.aliases, model.Alias{Name: "zen", Alias: "Zen Browser"})
	a.aliases = append(a.aliases, model.Alias{Name: "org.gnome.TextEditor", Alias: "Gnome Text Editor"})
	a.aliases = append(a.aliases, model.Alias{Name: "org.gnome.Nautilus", Alias: "Nautilus"})

	return a
}

func (a *alias) ReplaceAppId2Alias(r model.Report) model.Report {

	for _, alias := range a.aliases {
		if strings.Contains(r.Name, alias.Name) {
			r.Name = strings.ReplaceAll(r.Name, alias.Name, alias.Alias)
		}
	}

	return r
}
