# http_log_monitor
Go CLI application to watch w3c-format log files

This application will monitor a w3c-format log file, giving you interesting statistics and alerting when
the average number of hits for a two minute period is passed.

## To Install
`go get github.com/Aevin1387/http_log_monitor`

## Usage
`http_log_monitor -log='/path/to/log' -format='log format (e.x. %h %l %u %t "%r" %>s %b)' -alert-on=26`

## Options
-log: The path to the log file to monitor

-format: The format the log file is written to

-alert-on: The average number of hits within a two minute period to create an alert on.
