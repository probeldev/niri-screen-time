# Niri Screen Time

A utility that collects information about how much time you spend in each application.

![Report Example](https://github.com/probeldev/niri-screen-time/blob/main/screenshots/report.png?raw=true)


## Support Wayland compositor:

Niri

Hyprland

## Установка

### go 

    go install github.com/probeldev/niri-screen-time@latest     


### nix 
    
    nix profile install github:probeldev/niri-screen-time

## Usage 

Start the daemon:

    niri-screen-time -daemon 

View data (default: today's data):

    niri-screen-time 

View data from a selected date until today:
    
    niri-screen-time -from=2023-10-20

View data within a date range:

    niri-screen-time -from=2023-10-01 -to=2023-10-31 


## Subroutine and Website Configuration

To track specific websites or console commands within applications, use the configuration file:

**Location:**  
`~/.config/niri-screen-time/subprograms.json`

### Configuration Format
```json
[
  {
    "app_id": "application_identifier",
    "title": "partial_window_title",
    "alias": "display_name"
  }
]

```

### Configuration Examples

For terminal commands:

```json
{
  "app_id": "com.mitchellh.ghostty",
  "title": "nvim",
  "alias": "NeoVim Editor"
}
```

For websites in browser:

```json
{
  "app_id": "org.mozilla.firefox",
  "title": "GitHub",
  "alias": "GitHub"
}
```

### Правила работы:

The system checks if the title is contained in the window title

Case insensitive

When matched: adds the specified alias in parentheses to the app ID

If the app is in the list but the title doesn't match: adds (Other)

The config file is automatically created with examples on first launch

## Application Alias Configuration

To display custom application names in reports, you can configure aliases in the config file.
The system checks for partial matches (if the name is contained in the application title).

### Configuration File

**Location:**  

    ~/.config/niri-screen-time/alias.json

### Configuration Format

```json
[
  {
    "name": "original_application_name",
    "alias": "custom_display_name"
  },
  {
    "name": "org.telegram.desktop",
    "alias": "Telegram"
  }
]
```

Made in Belarus with ❤️
