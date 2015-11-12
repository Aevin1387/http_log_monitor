package main

import (
  "flag"
  "fmt"
  "io/ioutil"
  "runtime/debug"
)

var (
  logFile = flag.String("log", "/tmp/log", "Log File to watch")
  format = flag.String("format", "%h %l %u %t \"%r\" %>s %b", "Log Format")
)

func main() {
  flag.Parse()

  defer func() {
    if e := recover(); e != nil {
      trace := fmt.Sprintf("%s: %s", e, debug.Stack());
      ioutil.WriteFile("trace.txt", []byte(trace), 0644);
    }
  }()

  rawLineChan := make(chan string)
  parsedLineChan := make(chan LineData)

  go startWatch(*logFile, rawLineChan)
  go startStats(parsedLineChan)
  startParser(rawLineChan, parsedLineChan, *format)
}
