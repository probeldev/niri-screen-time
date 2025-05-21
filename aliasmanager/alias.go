package aliasmanager

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/probeldev/niri-screen-time/model"
)

type AliasManager struct {
	aliases    []model.Alias
	configPath string
}

func NewAliasManager() (*AliasManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".config", "niri-screen-time", "alias.json")

	am := &AliasManager{
		configPath: configPath,
	}

	if err := am.loadAliases(); err != nil {
		return nil, err
	}

	return am, nil
}

func (am *AliasManager) loadAliases() error {
	file, err := os.ReadFile(am.configPath)
	if err != nil {
		// Если файла нет, используем дефолтные значения
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if err := json.Unmarshal(file, &am.aliases); err != nil {
		return err
	}

	return nil
}

func (am *AliasManager) ReplaceAppId2Alias(r model.Report) model.Report {
	for _, alias := range am.aliases {
		if strings.Contains(r.Name, alias.Name) {
			r.Name = strings.ReplaceAll(r.Name, alias.Name, alias.Alias)
		}
	}
	return r
}
