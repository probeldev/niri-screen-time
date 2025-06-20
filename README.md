# Niri Screen Time

Утилита, которая собирает информацию о том, сколько времени проводишь в каждом из приложений.

![Report Example](https://github.com/probeldev/niri-screen-time/blob/main/screenshots/report.png?raw=true)


## Support Wayland compositor:

Niri

Hyprland

## Установка

### go 

    go install github.com/probeldev/niri-screen-time@latest     


### nix 
    
    nix profile install github:probeldev/niri-screen-time

## Использование 

Запуск демона

    niri-screen-time -daemon 

Просмотр данных(По умолчанию за сегодня)

    niri-screen-time 

Просмотр данных, от выбранной даты по сегодня 
    
    niri-screen-time -from=2023-10-20

Просмотр данных, диапазон данных

    niri-screen-time -from=2023-10-01 -to=2023-10-31 


## Конфигурация подпрограмм и сайтов

Для определения конкретных сайтов или консольных команд внутри приложений используется файл конфигурации:

**Расположение:**  
`~/.config/niri-screen-time/subprograms.json`

### Формат конфигурации
```json
[
  {
    "app_id": "идентификатор_приложения",
    "title": "часть_заголовка_окна",
    "alias": "отображаемое_название"
  }
]

```

### Примеры конфигурации

Для терминальных команд:

```json
{
  "app_id": "com.mitchellh.ghostty",
  "title": "nvim",
  "alias": "NeoVim Editor"
}
```

Для веб-сайтов в браузере:

```json
{
  "app_id": "org.mozilla.firefox",
  "title": "GitHub",
  "alias": "GitHub"
}
```

### Правила работы:

Система проверяет вхождение title в заголовок окна

Регистр символов не учитывается

При совпадении к app_id добавляется указанный alias в скобках

Если приложение есть в списке, но заголовок не совпал - добавляется "(Other)"

Файл автоматически создаётся с примерами при первом запуске

## Конфигурация алиасов приложений

Для удобного отображения названий приложений в отчётах вы можете настроить алиасы (псевдонимы) в файле конфигурации.
Система проверяет частичное совпадение (если name содержится в названии приложения).

### Файл конфигурации

Расположение:

    ~/.config/niri-screen-time/alias.json

### Формат файла

```json
[
  {
    "name": "оригинальное_название_приложения",
    "alias": "желаемое_название"
  },
  {
    "name": "org.telegram.desktop",
    "alias": "Telegram"
  }
]
```

Made in Belarus with ❤️
