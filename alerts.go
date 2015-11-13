package main

import (
  "fmt"
  "time"

  "github.com/nsf/termbox-go"
)

const (
  alertHistoryAmount = 12 // 2 Minutes
  alertLine = 5
)

var openAlert      *Alert           // an open alert, if any
var alerts         []*Alert         // all alerts
var startTime      time.Time

type Alert struct {
  hits          int               // number of hits causing alert
  alertedAt     time.Time         // time alert started
  resolvedAt    time.Time         // time alert resolved
  resolved      bool
}

func setupAlerts() {
  alerts = make([]*Alert, 0)
  startTime = time.Now()
}

func displayAlerts() {
  checkForAlerts()
  alertsToDisplay := alerts
  maxHeight := screenHeight - alertLine - 1
  totalAlerts := len(alertsToDisplay)

  if(totalAlerts > maxHeight) {
    alertsToDisplay = alertsToDisplay[totalAlerts - maxHeight:]
  }

  lineNum := alertLine
  for _, alert := range alertsToDisplay {
    line := fmt.Sprintf("Alert: %v hits on average at %v", alert.hits, alert.alertedAt)
    printAlert(0, lineNum, line)
    lineNum++
    if(alert.resolved) {
      line := fmt.Sprintf("Alert resolved at %v", alert.resolvedAt)
      printResolved(0, lineNum, line)
      lineNum++
    }
  }
}

func checkForAlerts() {
  hits := averageHitCount()
  if (hits >= *alertAmount && openAlert == nil) {
    openAlert = new(Alert)
    openAlert.hits = hits
    openAlert.alertedAt = time.Now()
    alerts = append(alerts, openAlert)
  } else if (canResolve()) {
    openAlert.resolvedAt = time.Now()
    openAlert.resolved = true
    openAlert = nil
  }
}

func averageHitCount() int {
  if(len(allSectionData) == 0) { return 0 }

  historyAmount := 0
  hits := 0
  for _, sectionData := range allSectionData {
    historyNeeded := sectionData.statHistory
    if len(historyNeeded) > alertHistoryAmount {
      index := len(historyNeeded) - alertHistoryAmount
      historyNeeded = sectionData.statHistory[index:]
    }

    if len(historyNeeded) > historyAmount { historyAmount = len(historyNeeded) }

    for _, stat := range historyNeeded {
      hits += stat.hitCount
    }
  }

  return (hits / historyAmount)
}

func canResolve() bool {
  // return openAlert != nil && time.Now().Sub(startTime) >= (time.Minute * 2)
  return openAlert != nil && time.Now().Sub(startTime) >= (time.Second * 30)
}

func printAlert(x, y int, msg string) {
  for _, c := range msg {
    termbox.SetCell(x, y, c, termbox.ColorRed, backgroundDefault)
    x++
  }
}

func printResolved(x, y int, msg string) {
  for _, c := range msg {
    termbox.SetCell(x, y, c, termbox.ColorGreen, backgroundDefault)
    x++
  }
}
