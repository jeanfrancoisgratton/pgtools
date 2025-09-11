// pgtools
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/07/09 10:58
// Original filename: src/db/logging.go

package logging

import (
	"log"
	"os"
	"path/filepath"
	"pgtools/types"
	"time"
)

var Logger *log.Logger

func Errorf(format string, v ...interface{}) {
	if types.LogLevel == "error" || types.LogLevel == "info" || types.LogLevel == "debug" {
		timestamp := time.Now().Format("2006/01/02 15:04:05")
		Logger.Printf("[ERROR] %s "+format, append([]interface{}{timestamp}, v...)...)
	}
}

func Infof(format string, v ...interface{}) {
	if types.LogLevel == "info" || types.LogLevel == "debug" {
		timestamp := time.Now().Format("2006/01/02 15:04:05")
		Logger.Printf("[INFO] %s "+format, append([]interface{}{timestamp}, v...)...)
	}
}

func Debugf(format string, v ...interface{}) {
	if types.LogLevel == "debug" {
		timestamp := time.Now().Format("2006/01/02 15:04:05")
		Logger.Printf("[DEBUG] %s "+format, append([]interface{}{timestamp}, v...)...)
	}
}

func init() {
	// Initialize log file
	base := filepath.Join(os.Getenv("HOME"), ".local", "state")
	f, _ := os.OpenFile(filepath.Join(base, "pgtools.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "WARNING: failed to open log file: %v\n", err)
	//	Logger = log.New(os.Stderr, "", log.LstdFlags)
	//} else {
	//	Logger = log.New(f, "", log.LstdFlags)
	//}
	Logger = log.New(f, "", log.LstdFlags)
}
