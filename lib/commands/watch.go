package commands

import (
	"github.com/thejezzi/watchmego/lib/cli"
	"github.com/thejezzi/watchmego/lib/logger"
	"github.com/thejezzi/watchmego/lib/watcher"
	"os"
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
