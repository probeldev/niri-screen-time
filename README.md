# Niri Screen Time
A utility that tracks how much time was spent on applications by their class and its title.

![Report Example](https://github.com/probeldev/niri-screen-time/blob/main/screenshots/report.png?raw=true)
## Supported Wayland Compositors
- Niri
- Hyprland
## Installation
### go
```bash
go install github.com/probeldev/niri-screen-time@latest
```
> [!NOTE]  
> If you get an error claiming that `niri-screen-time` cannot be found or that it is not defined, you
> may need to add `~/go/bin` to your $[PATH](https://jvns.ca/blog/2025/02/13/how-to-add-a-directory-to-your-path/).

### nix 

```bash 
nix profile install github:probeldev/niri-screen-time
```
> [!NOTE]
> This is an imperative way of installing, and it is not preferred on Nix.
> You should define it in your Nix configuration instead.
## Usage

### Daemon

Daemon — a background service that collects data on application usage time (requires adding to startup for proper operation).

Start the daemon:

```bash
niri-screen-time -daemon 
```

### Report 

(without key)

```bash
niri-screen-time 
```

View data from a selected date until today:
  
```bash
niri-screen-time -from=2023-10-20
```

View data within a date range:

```bash
niri-screen-time -from=2023-10-01 -to=2023-10-31 

```


#### Subroutine and Website Configuration

To track specific websites or console commands within applications, use the configuration file:

**Location:**  
`~/.config/niri-screen-time/subprograms.json`

##### Configuration Format
```json
[
  {
    "app_id": "application_identifier",
    "title": "partial_window_title",
    "alias": "display_name"
  }
]

```

##### Configuration Examples

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

##### How It Works:

The system checks if the title is contained in the window title

Case insensitive

When matched: adds the specified alias in parentheses to the app ID

If the app is in the list but the title doesn't match: adds (Other)

The config file is automatically created with examples on first launch

#### Application Alias Configuration

To display custom application names in reports, you can configure aliases in the config file.
The system checks for partial matches (if the name is contained in the application title).

##### Configuration File

**Location:**  

```
~/.config/niri-screen-time/alias.json
```

##### Configuration Format

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

### Details

This mod adds detailed per-application stats.

```bash
niri-screen-time -details -appid="org.telegram.desktop" -title="" -from='2025-01-20' -to='2025-08-20' -limit=20 -onlytex
t

```

## License  
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
