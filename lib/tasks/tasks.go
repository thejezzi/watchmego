package tasks

import (
	"github.com/thejezzi/watchmego/lib/checker"
	"github.com/thejezzi/watchmego/lib/cli"
	"github.com/thejezzi/watchmego/lib/logger"
	"os"
	"os/exec"
)

func RunMakeWatch(args *cli.Args) {

	logger.Info("Running make watch")
	checker.CheckMakefile(args.Dir)
	checker.CheckMakefileWatch(args.Dir)

	// run command "make watch"
	cmd := exec.Command("make", "watch")

	// set working directory
	cmd.Dir = args.Dir

	output, err := cmd.Output()

	if err != nil {
		logger.Error("Error running make watch")
		if args.Debug {
			logger.Debug(err.Error())
		}
		os.Exit(1)
	}

	logger.Info("\t" + string(output))
}
