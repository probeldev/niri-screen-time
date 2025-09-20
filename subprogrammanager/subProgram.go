// Package subprogrammanager - implementation subprogram
package subprogrammanager

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/probeldev/niri-screen-time/model"
)

type SubProgramManager struct {
	programs   []model.SubProgram
	configPath string
}

func NewSubProgramManager() (*SubProgramManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".config", "niri-screen-time", "subprograms.json")

	spm := &SubProgramManager{
		configPath: configPath,
	}

	if err := spm.loadPrograms(); err != nil {
		return nil, err
	}

	return spm, nil
}

func (spm *SubProgramManager) loadPrograms() error {
	file, err := os.ReadFile(spm.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if err := json.Unmarshal(file, &spm.programs); err != nil {
		return err
	}

	return nil
}

func (spm *SubProgramManager) IsSetProgram(st model.ScreenTime) bool {
	for _, p := range spm.programs {
		for _, id := range p.AppIDs {
			if id == st.AppID {
				return true
			}
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
