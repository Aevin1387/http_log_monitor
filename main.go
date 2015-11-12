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

  out := make(chan string)

  go startWatch(*logFile, out)
  startParser(out, *format)
}
