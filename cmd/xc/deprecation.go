package main

import (
	"log"
	"os"

	"github.com/joerdav/xc/models"
)

func warnInteractive(tasks models.Tasks) {
	var interactiveTasks []string
	for _, task := range tasks {
		if task.Interactive {
			interactiveTasks = append(interactiveTasks, task.Name)
		}
	}
	if len(interactiveTasks) > 0 {
		logger := log.New(os.Stderr, "warning: xc: ", 0)
		if len(interactiveTasks) == 1 {
			logger.Printf(`Task "%s" is set to Interactive,
             which is a deprecated attribute that has no effect.
             You can safely remove the attribute.
             (see https://github.com/joerdav/xc/issues/127)`, interactiveTasks[0])
		} else {
			logger.Printf(`Tasks %v are set to Interactive,
             which is a deprecated attribute that has no effect.
             You can safely remove the attribute.
             (see https://github.com/joerdav/xc/issues/127)`, interactiveTasks)
		}
	}
}
