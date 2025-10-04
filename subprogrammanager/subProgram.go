// Package subprogrammanager - implementation subprogram
package subprogrammanager

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/probeldev/niri-screen-time/model"
	"gopkg.in/yaml.v3"
)

type SubProgramManager struct {
	programs  []model.SubProgram
	configDir string
}

func NewSubProgramManager() (*SubProgramManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(homeDir, ".config", "niri-screen-time")

	spm := &SubProgramManager{
		configDir: configDir,
	}

	if err := spm.loadPrograms(); err != nil {
		return nil, err
	}

	return spm, nil
}

func (spm *SubProgramManager) loadPrograms() error {
	// Try YAML files first, then fall back to JSON
	yamlExtensions := []string{".yaml", ".yml", ".json"}

	for _, ext := range yamlExtensions {
		configFile := filepath.Join(spm.configDir, "subprograms"+ext)
		file, err := os.ReadFile(configFile)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}

		if ext == ".json" {
			if err := json.Unmarshal(file, &spm.programs); err != nil {
				return err
			}
		} else {
			if err := yaml.Unmarshal(file, &spm.programs); err != nil {
				return err
			}
		}

		return nil
	}

	return nil
}

func (spm *SubProgramManager) IsSetProgram(st model.ScreenTime) bool {
	for _, p := range spm.programs {
		if slices.Contains(p.AppIDs, st.AppID) {
			return true
		}
	}
	return false
}

func (spm *SubProgramManager) GetSubProgram(st model.ScreenTime) model.ScreenTime {
	if !spm.IsSetProgram(st) {
		return st
	}

	for _, p := range spm.programs {
		isCompareID := false
		for _, id := range p.AppIDs {
			if id == st.AppID {
				isCompareID = true
			}
		}

		if !isCompareID {
			continue
		}

		if len(p.TitleList) == 0 {
			st.AppID = p.Alias
			return st
		}

		for _, title := range p.TitleList {
			if strings.Contains(st.Title, title) {
				st.AppID = p.Alias
				return st
			}
		}
	}

	st.AppID += " (Other)"
	return st
}
