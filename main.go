package main

import (
	"fmt"
	"os"

	"github.com/thejezzi/watchmego/lib/checker"
	"github.com/thejezzi/watchmego/lib/cli"
	"github.com/thejezzi/watchmego/lib/commands"
	"github.com/thejezzi/watchmego/lib/gen"
	"github.com/thejezzi/watchmego/lib/logger"
	"github.com/thejezzi/watchmego/lib/util"
)

const version string = "0.0.1"

func main() {
	args := cli.ParseArgs(os.Args[1:])

	switch {
	case args.Version:
		fmt.Println("wmg " + version)
	case args.Check:
		if checker.CheckMakefile(args.Dir) {
			logger.Info("Makefile found")
			if checker.CheckMakefileWatch(args.Dir) {
				logger.Info("Makefile contains watch target")
			} else {
				logger.Error("Makefile does not contain watch target")
			}
		} else {
			logger.Error("Makefile not found")
			answer := util.AskYesNoQuestion("Create Makefile?")
			if answer {
				gen.CreateMakeFile(args.Dir)
			}
			logger.Info("Please edit your makefile and add a watch target.")
			logger.Info("Exiting ...")
			os.Exit(0)
		}
	case args.Create:
		if checker.CheckMakefile(args.Dir) {
			logger.Error("Makefile already exists")
		} else {
			gen.CreateMakeFile(args.Dir)
		}
	default:
		logger.PrintGopher()
		commands.Watch(&args)
	}
}
