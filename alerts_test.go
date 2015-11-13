package main

import (
  "fmt"
  "testing"
  "time"
)


func TestCanResolve(t *testing.T) {
  startTime = time.Now().Add(-time.Minute * 2)
  if canResolve() {
    t.Error("canResolve alerted that it can resolve when there is nothing to resolve")
  }

  openAlert = new(Alert)
  if !canResolve() {
    t.Error("canResolve alerted that it can resolve when it couldn't")
  }

  startTime = time.Now().Add(-time.Second)
  if canResolve() {
    t.Error("canResolve alerted that it can't resolve when it should")
  }
}

func TestAverageHitCount(t *testing.T) {
  if(averageHitCount() != 0) {
    t.Error("averageHitCount returned more than 0 when there are no hits")
  }

  allSectionData = make(map[string]*SectionData, 0)
  sectionData := dataForSection("test")
  sectionData.currentStat.hitCount = *alertAmount

  if(averageHitCount() != *alertAmount) {
     t.Error("averageHitCount did not return the amount from a single stat")
  }

  sectionData.StartNewStats()
  sectionData.currentStat.hitCount = *alertAmount * 3
  if(averageHitCount() != *alertAmount * 2) {
     t.Error("averageHitCount did not return the average amount from stat history")
  }

  for i := 0; i < alertHistoryAmount; i++ {
    sectionData.StartNewStats()
    sectionData.currentStat.hitCount  = 1
  }

  if(averageHitCount() != 1) {
     t.Error("averageHitCount returned ", averageHitCount(), " instead of 1 (the average)")
  }

  newSectionData := dataForSection("test2")
  for i := 0; i < alertHistoryAmount; i++ {
    newSectionData.StartNewStats()
    newSectionData.currentStat.hitCount  = 1
  }
  if(averageHitCount() != 2) {
     t.Error("averageHitCount did not include all sections, amount was", averageHitCount(), "expected 2")
  }
}

func TestCheckForAlerts(t *testing.T) {
  allSectionData = make(map[string]*SectionData, 0)
  sectionData := dataForSection("test")
  sectionData.currentStat.hitCount = *alertAmount
  openAlert = nil

  checkForAlerts()
  if(len(alerts) != 1) {
    t.Error("checkForAlerts did not open an alert")
  }

  openAlert = new(Alert)
  alerts = make([]*Alert, 0)
  checkForAlerts()
  if(len(alerts) != 0) {
    t.Error("checkForAlerts opened an alert when there was one")
  }

  startTime = time.Now().Add(-time.Minute * 2)
  alerts = append(alerts, openAlert)
  checkForAlerts()
  if(openAlert != nil) {
    t.Error("checkForAlerts did not resolve open alert")
  }

  alert := alerts[0]
  if(!alert.resolved) {
    t.Error("checkForAlerts did not mark the alert resolved")
  }

  if(alert.resolvedAt.IsZero()) {
    t.Error("checkForAlerts did not set the resolve time")
  }
}

func TestAlertsToDisplay(t *testing.T) {
  screenHeight = 6

  alertLen := len(alertsToDisplay())
  if(alertLen != 0) {
    t.Error("alertsToDisplay returned", alertLen, "should have been 0 due to height")
  }

  screenHeight = 7
  alert := new(Alert)
  alerts = append(alerts, alert)

  alertLen = len(alertsToDisplay())
  if(alertLen != 1) {
    t.Error("alertsToDisplay returned", alertLen, "should have been 1")
  }

  alert = new(Alert)
  alerts = append(alerts, alert)
  alertLen = len(alertsToDisplay())
  if(alertLen != 1) {
    t.Error("alertsToDisplay returned", alertLen, "should have been 1 due to height")
  }

  screenHeight = 8
  alertLen = len(alertsToDisplay())
  if(alertLen != 2) {
    t.Error("alertsToDisplay returned", alertLen, "should have been 2 due to height")
  }
}

// MockPrint is a container for data passed to the mock
// printing functions
type MockPrint struct {
  x int
  y int
  msg string
}

var mockAlertPrints []MockPrint
var mockResolvePrints []MockPrint
var mockAllPrints []MockPrint

func mockPrintAlert(x, y int, msg string) {
  mockPrint := MockPrint{ x: x, y: y, msg: msg}
  mockAlertPrints = append(mockAlertPrints, mockPrint)
  mockAllPrints = append(mockAllPrints, mockPrint)
}

func mockPrintResolved(x, y int, msg string) {
  mockPrint := MockPrint{ x: x, y: y, msg: msg}
  mockResolvePrints = append(mockResolvePrints, mockPrint)
  mockAllPrints = append(mockAllPrints, mockPrint)
}

func TestDisplayAlerts(t *testing.T) {
  screenHeight = 30
  alert := &Alert{ hits: 10, alertedAt: time.Now() }
  alerts = make([]*Alert, 0)
  alerts = append(alerts, alert)
  allSectionData = make(map[string]*SectionData, 0)

  displayAlerts(mockPrintAlert, mockPrintResolved)

  numAlerts := len(mockAlertPrints)
  numResolveds := len(mockResolvePrints)
  if(numAlerts != 1) {
    t.Error("displayAlerts did not print the alert")
  }

  if(numResolveds != 0) {
    t.Error("displayAlerts printed a resolved when the alert was not")
  }

  alert.resolved = true
  alert.resolvedAt = time.Now()
  mockAlertPrints = make([]MockPrint,0)
  mockResolvePrints = make([]MockPrint,0)

  displayAlerts(mockPrintAlert, mockPrintResolved)
  numAlerts = len(mockAlertPrints)
  numResolveds = len(mockResolvePrints)

  if(numAlerts != 1) {
    t.Error("displayAlerts did not print the alert")
  }

  if(numResolveds != 1) {
    t.Error("displayAlerts did not print a resolution when it should")
  }

  mockAlertPrints = make([]MockPrint,0)
  mockResolvePrints = make([]MockPrint,0)
  mockAllPrints = make([]MockPrint,0)
  alert = &Alert{ hits: 10, alertedAt: time.Now() }
  alerts = append(alerts, alert)
  displayAlerts(mockPrintAlert, mockPrintResolved)
  numAlerts = len(mockAlertPrints)
  numResolveds = len(mockResolvePrints)

  if(numAlerts != 2) {
    t.Error("displayAlerts did not print the right alerts")
  }

  if(numResolveds != 1) {
    t.Error("displayAlerts did not print a resolution when it should")
  }

  expectedLine := fmt.Sprintf("Alert: %v hits on average at %v", alerts[0].hits, alerts[0].alertedAt)
  if(mockAllPrints[0].msg != expectedLine) {
    t.Error("First alert had wrong message")
  }

  if(mockAllPrints[0].y != 5) {
    t.Error("First alert printed on wrong line", mockAllPrints[0].y)
  }

  expectedLine = fmt.Sprintf("Alert resolved at %v", alerts[0].resolvedAt)
  if(mockAllPrints[1].msg != expectedLine) {
    t.Error("First resolution had wrong message", mockAllPrints[1].msg)
  }

  if(mockAllPrints[1].y != 6) {
    t.Error("First resolution printed on wrong line", mockAllPrints[1].y)
  }

  expectedLine = fmt.Sprintf("Alert: %v hits on average at %v", alerts[1].hits, alerts[1].alertedAt)
  if(mockAllPrints[2].msg != expectedLine) {
    t.Error("Third alert had wrong message", mockAllPrints[2].msg)
  }

  if(mockAllPrints[2].y != 7) {
    t.Error("Third alert printed on wrong line", mockAllPrints[2].y)
  }
}
