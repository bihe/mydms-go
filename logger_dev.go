// +build !prod

package main

import (
	"log"
	"os"
)

// InitLogger performs a setup for the logging mechanism
func InitLogger(conf LogConfig) {
	log.SetPrefix(LogPrefix(conf))
	log.SetOutput(os.Stdout)
}
