package main

import (
    "os"
    log "github.com/Sirupsen/logrus"
)

func main() {
    initLogging()
	Control()
}

func initLogging() {
  log.SetOutput(os.Stderr)
  log.SetLevel(log.DebugLevel)
}

