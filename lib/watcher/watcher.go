package watcher

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/thejezzi/watchmego/lib/cli"
	"github.com/thejezzi/watchmego/lib/logger"
	"github.com/thejezzi/watchmego/lib/util"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	Watcher        *fsnotify.Watcher
	Args           cli.Args
	Callbacks      []func(*cli.Args, *chan bool)
	alreadyRunning bool
}

func New(args cli.Args) *Watcher {
	W := &Watcher{}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Error("Error creating watcher")
	}
	W.Watcher = watcher
	W.Args = args

	_, errDir := filepath.Abs(args.Dir)
	if errDir != nil {
		logger.Error("Error getting absolute path")
		if W.Args.Debug {
			logger.Debug(errDir.Error())
		}
		os.Exit(1)
	}

	dir := args.Dir

	filePathError := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Error("Error walking directory")
			if W.Args.Debug {
				logger.Debug(err.Error())
			}
			os.Exit(1)
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			logger.Info("Watching file: " + path)
			err = W.Watcher.Add(path)
			if err != nil && W.Args.Debug {
				logger.Error("Error adding file to watcher")
			}
		}
		return nil
	})

	if filePathError != nil {
		logger.Error("Error walking directory")
		if W.Args.Debug {
			logger.Debug(filePathError.Error())
		}
		os.Exit(1)
	}

	return W
}

func (W *Watcher) Close() {
	W.Watcher.Close()
}

func (W *Watcher) Watch() {
	// Run the build step
	initial_err := W.runBuildStep()

	// Run once command potentially endlessly
	go runOnce(&W.Args)

	stop := make(chan bool)

	if initial_err != nil {
		go runEmpty(&stop, &W.Args)
	} else {
		go runServer(&stop, &W.Args)
	}

	for {
		select {
		case event, ok := <-W.Watcher.Events:
			if !ok {
				if W.Args.Debug {
					logger.Debug("Watcher closed")
				}
				return
			}
			if W.Args.Debug {
				logger.Debug("Event: " + event.String())
				logger.Debug("Event name " + event.Name)
				logger.Debug("File extension: " + filepath.Ext(event.Name))
				logger.Debug("Event has write " + strconv.FormatBool(event.Has(fsnotify.Chmod)))
				fmt.Println()
			}
			// Rerun makescript if file is modified in any way
			if !W.alreadyRunning && filepath.Ext(event.Name) == ".go" && (event.Has(fsnotify.Chmod) || event.Has(fsnotify.Write)) {
				W.alreadyRunning = true
				logger.Info("File modified: " + event.Name + ", rerunning makescript")

				stop <- true
				err := W.runBuildStep()
				stop = make(chan bool)

				if err != nil {
					go runEmpty(&stop, &W.Args)
					continue
				}

				go runServer(&stop, &W.Args)

			} else {
				W.alreadyRunning = false
			}

			// Remove and re-add file to watcher because iotify doesn't work properly

			// err := W.Watcher.Remove(event.Name)
			// if err != nil {
			//   logger.Error("Error removing file from watcher" + err.Error())
			//   continue
			// }

			err := W.Watcher.Add(event.Name)
			if err != nil && W.Args.Debug {
				logger.Error("Error adding file to watcher")
			}

		case err, ok := <-W.Watcher.Errors:
			if !ok {
				if W.Args.Debug {
					logger.Debug("Watcher closed")
				}
				return
			}
			logger.Error("Error: " + err.Error())
			// Stop the watcher
		}
	}
}

func (W *Watcher) runBuildStep() error {
	cmd := exec.Command("make", "watch-build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		if W.Args.Debug {
			logger.Error("Error running make watch")
			logger.Debug(err.Error())
		}
		return err
	}

	// Wait for the process to finish
	err = cmd.Wait()
	if err != nil {
		if W.Args.Debug {
			logger.Error("Error running make watch")
			logger.Debug(err.Error())
		}
		return err
	}

	return nil
}

func runOnce(args *cli.Args) {
	// Create signal which never ends
	infinity := make(chan bool)

	once := util.ReadMakeFileVar("WATCHMEGO_ONCE")

	if once == "" {
		return
	}

	logger.Info("Running command: '" + once + "' once (from Makefile)")

	commandParts := strings.Split(strings.TrimSpace(once), " ")
	commandName := commandParts[0]
	commandRest := strings.Join(commandParts[1:], " ")
	cmd := exec.Command(commandName, util.ParseArguments(commandRest)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		logger.Error("Error running command")
		if args.Debug {
			logger.Debug(err.Error())
		}
	}

	for {
		select {
		case <-infinity:
			os.Exit(0)
		}
	}
}

func runServer(stop *chan bool, args *cli.Args) {
	// Get the run command from makefile
	runCommand := util.ReadMakeFileVar("WATCHMEGO")
	allCommands := strings.Split(runCommand, ";")

	var allCommandHandles []*exec.Cmd

	if allCommands[0] == "" {
		logger.Error("No run command found in Makefile. Consider defining WATCHMEGO in your Makefile")
		// os.Exit(1)
	}

	for _, command := range allCommands {
		if command == "" {
			continue
		}
		logger.Info("Running command: '" + command + "' (from Makefile)")
		// Get first word of command
		commandParts := strings.Split(strings.TrimSpace(command), " ")
		commandName := commandParts[0]
		commandRest := strings.Join(commandParts[1:], " ")
		cmd := exec.Command(commandName, util.ParseArguments(commandRest)...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			logger.Error("Error running command")
			if args.Debug {
				logger.Debug(err.Error())
			}
		}
		allCommandHandles = append(allCommandHandles, cmd)
	}

	// Run the command
	// cmd := exec.Command(strings.TrimSpace(runCommand))
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	//
	// err := cmd.Start()
	// if err != nil {
	// 	logger.Error("Error running command")
	// 	if args.Debug {
	// 		logger.Debug(err.Error())
	// 	}
	// }

	<-*stop

	// // Kill the process
	// err = cmd.Process.Kill()
	// if err != nil {
	// 	if args.Debug {
	// 		logger.Error("Error killing process")
	// 		logger.Debug(err.Error())
	// 	}
	// }
	//
	// // Wait for the process to finish
	// err = cmd.Wait()
	// if err != nil && args.Debug {
	// 	logger.Debug("Process killed")
	// }

	for _, cmd := range allCommandHandles {
		// Kill the process
		err := cmd.Process.Kill()
		if err != nil {
			if args.Debug {
				logger.Error("Error killing process")
				logger.Debug(err.Error())
			}
		}

		// Wait for the process to finish
		err = cmd.Wait()
		if err != nil && args.Debug {
			logger.Debug("Process killed")
		}
	}
}

func runEmpty(stop *chan bool, args *cli.Args) {
	<-*stop
	if args.Debug {
		logger.Debug("Process killed")
	}
}
