package main

import (
	"fmt"

	"github.com/bihe/mydms/core"
)

// LogPrefix is used to display a meaningful prefix for log-messages
func LogPrefix(config core.LogConfig) string {
	return fmt.Sprintf("[%s] ", config.Prefix)
}
