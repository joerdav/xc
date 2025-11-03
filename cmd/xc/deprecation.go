package main

import (
	"log"
	"os"

	"github.com/joerdav/xc/models"
)

func warnInteractive(tasks models.Tasks) {
	var logger = log.New(os.Stderr, "warning: xc: ", 0)
	for _, task := range tasks {
		if task.Interactive {
			logger.Printf(`Task "%s" is set to Interactive,
             which is a deprecated attribute that has no effect.
             You can safely remove the attribute.
             (see https://github.com/joerdav/xc/issues/127)`, task.Name)
		}
	}
}
