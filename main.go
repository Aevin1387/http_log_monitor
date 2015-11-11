package main

import (
  "flag"
  "fmt"
  "io/ioutil"
  "runtime/debug"
)

var (
  logFile = flag.String("log", "/tmp/log", "Log File to watch")
)

func main() {
  flag.Parse()

  defer func() {
    if e := recover(); e != nil {
      trace := fmt.Sprintf("%s: %s", e, debug.Stack());
      ioutil.WriteFile("trace.txt", []byte(trace), 0644);
    }
  }()
}