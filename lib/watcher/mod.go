package watcher

import (
  "fmt"
  "os"
  "path/filepath"
  "strconv"
  "wmg/lib/cli"
  "wmg/lib/logger"

  "github.com/fsnotify/fsnotify"
)

type Watcher struct {
  Watcher *fsnotify.Watcher
  Args cli.Args
  Callbacks []func(*cli.Args)
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

  err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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
      if err != nil {
        logger.Error("Error adding file to watcher")
      }
    }
    return nil
  })

  return W
}

func (W *Watcher) Close() {
  W.Watcher.Close()
}

func (W *Watcher) Watch() {
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
      if filepath.Ext(event.Name) == ".go" && (event.Has(fsnotify.Chmod) || event.Has(fsnotify.Write)){
        logger.Info("File modified: " + event.Name + ", rerunning makescript")
        W.RunCallbacks()
      }
      // Remove and re-add file to watcher because iotify doesn't work properly
      // err := W.Watcher.Remove(event.Name)
      // if err != nil {
      //   logger.Error("Error removing file from watcher" + err.Error())
      // }
      err := W.Watcher.Add(event.Name)
      if err != nil {
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
    }
  }
}

func (W *Watcher) AddCallback(callback func(*cli.Args)) {
  W.Callbacks = append(W.Callbacks, callback)
}

func (W *Watcher) RunCallbacks() {
  for _, callback := range W.Callbacks {
    callback(&W.Args)
  }
}



