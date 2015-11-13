package main

import (
  "fmt"
  "io/ioutil"
  "runtime/debug"
  "sort"
  "strings"
  // "strconv"
  "time"

  "github.com/nsf/termbox-go"
)

var quit = false
var screenWidth, screenHeight int
var statStart time.Time

const (
  backgroundDefault = termbox.ColorDefault
  statDisplayTime = time.Second * 10
)

func startDisplay() {
  defer func() {
    if e := recover(); e != nil {
      termbox.Close()
      trace := fmt.Sprintf("%s: %s", e, debug.Stack());
      ioutil.WriteFile("trace.txt", []byte(trace), 0644);
      quit = true
    }
  }()

  err := termbox.Init()
  if err != nil {
    panic(err)
  }
  termbox.SetInputMode(termbox.InputEsc)

  setupEvents()
  setupAlerts()
  printLineTo(0, 0, termbox.ColorWhite, backgroundDefault, "Waiting for hits")

  timer := time.Tick(time.Millisecond * 100)
  statStart = time.Now()
  for {
    select {
    case <-timer:
      redrawScreen()
    }
  }
}

func setupEvents() {
  eventChan := make(chan termbox.Event, 16)
  go handleEvents(eventChan)
  go func() {
    for {
      ev := termbox.PollEvent()
      eventChan <- ev
    }
  }()
}

func handleEvents(eventChan chan termbox.Event) {
  for {
    ev := <-eventChan
    switch ev.Type {
    case termbox.EventKey:
      switch ev.Key {
      case termbox.KeyEsc:
        fallthrough
      case termbox.KeyCtrlC:
        fallthrough
      case termbox.KeyCtrlQ:
         quit = true
      }
    case termbox.EventError:
      panic(ev.Err)
    }
  }
}

func redrawScreen() {
  screenWidth, screenHeight = termbox.Size()

  if(time.Now().Sub(statStart) > statDisplayTime) {
    termbox.Clear(backgroundDefault, backgroundDefault)
    displayStatistics()
    displayAlerts(printAlert, printResolved)
    statStart = time.Now()
    archiveStatData()
  }

  termbox.HideCursor()
  termbox.Flush()
}

func printLineTo(x, y int, fg, bg termbox.Attribute, msg string) {
  for _, c := range msg {
    termbox.SetCell(x, y, c, fg, bg)
    x++
  }
}

func printNormalLine(x, y int, msg string) {
  for _, c := range msg {
    termbox.SetCell(x, y, c, termbox.ColorWhite, backgroundDefault)
    x++
  }
}

func displayStatistics() {
  hits := 0
  methodCounts := make(map[string]int, 0)
  statusCounts := make(map[string]int, 0)
  hostCounts := make(map[string]int, 0)
  sectionHits := make(map[string]int, 0)

  for _, sectionData := range allSectionData {
    for _, stat := range sectionData.statHistory {
      hits += stat.hitCount
      sectionHits[sectionData.name] += stat.hitCount
      sumCounts(methodCounts, stat.methodsCounts)
      sumCounts(statusCounts, stat.statusCounts)
      sumCounts(hostCounts, stat.hostCounts)
    }
  }

  statsLine := fmt.Sprintf("Total Hits: %v. ", hits)
  printLineTo(0, 0, termbox.ColorWhite, backgroundDefault, statsLine)

  if(hits > 0) {
    printLineTo(0, 1, termbox.ColorWhite, backgroundDefault, " Top Sections: " + topThreeHits(sectionHits))
    printLineTo(0, 2, termbox.ColorWhite, backgroundDefault, " Top Methods: " + topThreeHits(methodCounts))
    printLineTo(0, 3, termbox.ColorWhite, backgroundDefault, " Top Statuses: " + topThreeHits(statusCounts))
    printLineTo(0, 4, termbox.ColorWhite, backgroundDefault, " Top Remote Hosts: " + topThreeHits(hostCounts))
  }
}

func sumCounts(totalCounts map[string]int, sectionCounts map[string] int) {
  for name, count := range sectionCounts {
    totalCounts[name] += count
  }
}

func topThreeHits(mapping map[string]int) string {
  stringCounts := mappingToStringCounts(mapping)
  sort.Sort(sort.Reverse(stringCounts))

  count := len(stringCounts)
  var arr []StringCount
  if(count < 3) {
    arr = stringCounts
  } else {
    arr = stringCounts[:3]
  }

  var countStrings []string
  for _, stringCount := range arr {
    countStrings = append(countStrings, stringCount.ToPrintString())
  }

  return strings.Join(countStrings, ", ")
}
