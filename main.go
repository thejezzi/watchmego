package main

import (
	"os"
	"wmg/lib/checker"
	"wmg/lib/cli"
	"wmg/lib/commands"
	"wmg/lib/gen"
	"wmg/lib/logger"
	"wmg/lib/util"
)

func main() {
  logger.PrintGopher()

  args := cli.ParseArgs(os.Args[1:])

  switch {
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
    commands.Watch(&args)
  }

}


