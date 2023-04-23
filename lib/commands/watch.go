package commands

import (
	"os"
	"wmg/lib/cli"
	"wmg/lib/logger"
	"wmg/lib/watcher"
)

func Watch(args *cli.Args) {
  w := watcher.New(*args) 
  // Dont forget to close the watcher
  defer w.Close()
  
  go w.Watch()

  <-make(chan struct{})

  logger.Info("Exiting...")
  os.Exit(0)
}
