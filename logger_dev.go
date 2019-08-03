// +build !prod

package main

import (
	"log"
	"os"

	"github.com/bihe/mydms/core"
)

// InitLogger performs a setup for the logging mechanism
func InitLogger(conf core.LogConfig) {
	log.SetPrefix(LogPrefix(conf))
	log.SetOutput(os.Stdout)
}
