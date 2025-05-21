# Niri Screen Time

Утилита, которая собирает информацию о том, сколько времени проводишь в каждом из приложений, для Wayland композитора Niri.

![Report Example](https://github.com/probeldev/niri-screen-time/blob/main/screenshots/report.png?raw=true)


## Установка

TODO

## Использование 

Запуск демона

    niri-screen-time -daemon 

Просмотр данных(По умолчанию за сегодня)

    niri-screen-time 

Просмотр данных, от выбранной даты по сегодня 
    
    niri-screen-time -from=2023-10-20

Просмотр данных, диапазон данных

    niri-screen-time -from=2023-10-01 -to=2023-10-31 
