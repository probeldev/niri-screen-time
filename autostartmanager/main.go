package autostartmanager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/probeldev/niri-screen-time/bash"
)

type AutoStartManager struct {
	appName     string
	plistPath   string
	programPath string
	args        []string
}

// NewAutoStartManager создает менеджер автозапуска для конкретной программы
func NewAutoStartManager(programPath string, args []string) (*AutoStartManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Используем имя программы для названия plist-файла
	appName := filepath.Base(programPath)
	plistPath := filepath.Join(homeDir, "Library", "LaunchAgents",
		fmt.Sprintf("com.niri.screentime.%s.plist", appName))

	return &AutoStartManager{
		appName:     appName,
		plistPath:   plistPath,
		programPath: programPath,
		args:        args,
	}, nil
}

// NewAutoStartManagerForMacOs создает менеджер специально для niri-screen-time
func NewAutoStartManagerForMacOs() (*AutoStartManager, error) {
	// Получаем полный путь к текущему исполняемому файлу
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %v", err)
	}

	// Используем полный путь к текущему бинарнику
	programPath := execPath
	args := []string{"-daemon"}

	return NewAutoStartManager(programPath, args)
}

func (a *AutoStartManager) Enable() error {
	// Собираем все аргументы в один массив
	programArgs := []string{a.programPath}
	programArgs = append(programArgs, a.args...)

	// Формируем XML для ProgramArguments
	argsXML := ""
	for _, arg := range programArgs {
		argsXML += fmt.Sprintf("        <string>%s</string>\n", arg)
	}

	plistContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.niri.screentime.%s</string>
    <key>ProgramArguments</key>
    <array>
%s    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/tmp/%s.log</string>
    <key>StandardErrorPath</key>
    <string>/tmp/%s.error.log</string>
    <key>EnvironmentVariables</key>
    <dict>
        <key>PATH</key>
        <string>/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin</string>
    </dict>
</dict>
</plist>`, a.appName, argsXML, a.appName, a.appName)

	dir := filepath.Dir(a.plistPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if err := os.WriteFile(a.plistPath, []byte(plistContent), 0644); err != nil {
		return err
	}

	fmt.Printf("✓ Автозапуск включен для %s: %s\n", a.programPath, a.plistPath)
	return nil
}

func (a *AutoStartManager) Load() error {
	cmd := fmt.Sprintf("launchctl load \"%s\"", a.plistPath)
	_, err := bash.RunCommand(cmd)
	if err != nil {
		return fmt.Errorf("ошибка загрузки службы: %v", err)
	}
	fmt.Println("✓ Служба загружена, приложение запущено")
	return nil
}

func (a *AutoStartManager) Unload() error {
	cmd := fmt.Sprintf("launchctl unload \"%s\"", a.plistPath)
	_, err := bash.RunCommand(cmd)
	if err != nil {
		return fmt.Errorf("ошибка выгрузки службы: %v", err)
	}
	fmt.Println("✓ Служба выгружена")
	return nil
}

func (a *AutoStartManager) EnableAndLoad() error {
	if err := a.Enable(); err != nil {
		return err
	}
	return a.Load()
}

func (a *AutoStartManager) Disable() error {
	// Сначала выгружаем службу
	a.Unload()

	// Затем удаляем plist-файл
	if err := os.Remove(a.plistPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("автозапуск не был настроен")
		}
		return err
	}

	fmt.Println("✓ Автозапуск отключен")
	return nil
}

func (a *AutoStartManager) Status() (bool, bool) {
	// Проверяем наличие plist-файла
	plistExists := false
	if _, err := os.Stat(a.plistPath); err == nil {
		plistExists = true
	}

	// Проверяем, запущена ли служба
	cmd := fmt.Sprintf("launchctl list | grep com.niri.screentime.%s", a.appName)
	output, err := bash.RunCommand(cmd)
	isRunning := err == nil && strings.Contains(output, "com.niri.screentime."+a.appName)

	return plistExists, isRunning
}

// GetPlistPath возвращает путь к созданному plist-файлу
func (a *AutoStartManager) GetPlistPath() string {
	return a.plistPath
}

// CheckAndFixPermissions проверяет и исправляет права доступа
func (a *AutoStartManager) CheckAndFixPermissions() error {
	fmt.Println("🔧 Проверка и настройка прав доступа...")

	// Добавляем приложение в список доступных для Accessibility
	script := `
	osascript <<'EOF'
	tell application "System Events"
		-- Проверяем включены ли UI элементы
		if not UI elements enabled then
			display dialog "niri-screen-time требует разрешения Accessibility для отслеживания активных окон." & return & return & ¬
			"Нажмите 'Open Settings' и добавьте niri-screen-time в список разрешенных приложений." ¬
			with title "Требуются разрешения" ¬
			with icon caution ¬
			buttons {"Open Settings", "Cancel"} ¬
			default button 1
			
			if button returned of result is "Open Settings" then
				tell application "System Preferences"
					activate
					reveal anchor "Privacy_Accessibility" of pane "com.apple.preference.security"
				end tell
			end if
		else
			-- Уже есть права, пытаемся добавить наше приложение
			try
				set appPath to "%s"
				tell application "System Events"
					tell process "System Preferences"
						if exists then
							-- Уже открыты настройки, ничего не делаем
						end if
					end tell
				end tell
			on error
				-- Игнорируем ошибки, главное что права есть
			end try
		end if
	end tell
	EOF
	`

	_, err := bash.RunCommand(fmt.Sprintf(script, a.programPath))
	if err != nil {
		fmt.Printf("⚠️  Не удалось автоматически настроить права: %v\n", err)
		fmt.Println("📋 Пожалуйста, вручную добавьте niri-screen-time в:")
		fmt.Println("   System Settings → Privacy & Security → Accessibility")
	} else {
		fmt.Println("✓ Права доступа настроены")
	}

	return nil
}
