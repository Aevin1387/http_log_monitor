package main

import (
  "flag"
  "fmt"
  "io/ioutil"
  "runtime/debug"
  "time"

  "github.com/nsf/termbox-go"
)

var (
  logFile = flag.String("log", "/tmp/log", "Log File to watch")
  format = flag.String("format", "%h %l %u %t \"%r\" %>s %b", "Log Format")
  alertAmount = flag.Int("alert-on", 26, "Amount of hits to alert on")
)

func main() {
  flag.Parse()

  defer func() {
    if e := recover(); e != nil {
      termbox.Close()
      trace := fmt.Sprintf("%s: %s", e, debug.Stack());
      ioutil.WriteFile("trace.txt", []byte(trace), 0644);
      quit = true
    }
  }()

  rawLineChan := make(chan string)
  parsedLineChan := make(chan LineData)

  go startDisplay()
  go startWatch(*logFile, rawLineChan)
  go startStats(parsedLineChan)

  go startParser(rawLineChan, parsedLineChan, *format)

  for {
    time.Sleep(time.Millisecond * 100)

    if quit == true {
      break
    }
  }

  termbox.Close()
}
