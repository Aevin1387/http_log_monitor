package main

import (
  "bufio"
  "log"
  "os"

  "gopkg.in/fsnotify.v1"
)

func startWatch(filename string, out chan string) {
  // Create Watcher
  watcher, err := fsnotify.NewWatcher()
  if err != nil {
    log.Fatal(err)
  }
  defer watcher.Close()

  file, err := os.Open(filename)
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  fileStats, err := os.Stat(filename)
  if err != nil {
    log.Fatal(err)
  }

  // Get current size of file
  lastSize := fileStats.Size()

  // Start Watcher
  err = watcher.Add(filename);
  if err != nil {
    log.Fatal(err)
  }

  // Watch file for writes
  done := make(chan bool)
  func() {
    for {
      select {
      case event := <-watcher.Events:
        if event.Op&fsnotify.Write == fsnotify.Write {
          file.Seek(0, 0) // Seek to beginning of file
          fileStats, err := os.Stat(filename)
          if err != nil {
            panic(err)
          }

          // Seek to previous end of file
          file.Seek(lastSize, 0)
          scanner := bufio.NewScanner(file)
          for scanner.Scan() {
            out <- scanner.Text() // Send line out
          }

          // Set new end of file
          lastSize = fileStats.Size()
        }
      case err := <-watcher.Errors:
        log.Println("error:", err);
      }
    }
  }()
  <- done
}
