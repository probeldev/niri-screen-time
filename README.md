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

### [aur](https://aur.archlinux.org/packages/niri-screen-time-git) (unofficial)

```bash
yay -S niri-screen-time-git
```

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
    "app_id": [
        "application_identi1",
        "application_identi2"
    ],
    "title": [
        "partial_window_title1",
        "partial_window_title2",
    ],
    "alias": "display_name"
  }
]

```

##### Configuration Examples

For terminal commands:

```json
{
    "app_ids": [
        "com.mitchellh.ghostty"
    ],
    "title_list": [
        "NeoVim: ~/script/",
        "NeoVim: ~/.config"
    ],
    "alias": "NeoVim: Редактирование ~/.config"
}
```

For websites in browser:

```json
{
    "app_ids": [
        "zen",
        "app.zen_browser.zen"
    ],
    "title_list":[
        "Monkeytype"
    ],
    "alias": "Monkeytype"
}
```

Just alias:
```json
{
    "app_ids": [
        "org.gnome.Nautilus"
    ],
    "title_list": [ ],
    "alias": "Nautilus"
}
```


### Details

This mod adds detailed per-application stats.

```bash
niri-screen-time -details -appid="org.telegram.desktop" -title="" -from='2025-01-20' -to='2025-08-20' -limit=20 -onlytex
t

```

## License  
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
