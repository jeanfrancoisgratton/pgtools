package main

import (
	"fmt"
	"os"
	"path/filepath"
	"pgtools/cmd"
	"time"
)

func main() {
	var err error
	var CurrentWorkingDir string

	// Create the pgtools logdir
	base := filepath.Join(os.Getenv("HOME"), ".local", "state")

	if err = os.MkdirAll(base, 0755); err != nil {
		fmt.Println(err)
		os.Exit(11)
	}

	// "Touch" the logfile so we avoid an error at the very first use of the tool
	if f, e := os.OpenFile(filepath.Join(base, "pgtools.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); e != nil {
		fmt.Println(e)
		os.Exit(12)
	} else {
		f.Close()
	}

	// Whatever happens, we need to preserve the current pwd, and restore it on exit, however the software exits
	if CurrentWorkingDir, err = os.Getwd(); err != nil {
		fmt.Println(err)
		os.Exit(13)
	}

	// We need to create a configuration directory. This is a per-user config dir
	if err = os.MkdirAll(filepath.Join(os.Getenv("HOME"), ".config", "JFG", "pgtools"), os.ModePerm); err != nil {
		fmt.Println(err)
		os.Exit(11)
	}

	// Mark execution boundaries in logfile : exec start
	if e := markLogFile("Start"); e != nil {
		fmt.Println(e)
		os.Exit(14)
	}

	// Command loop
	cmd.Execute()

	// Software execution is complete, let's get the hell outta here
	_ = os.Chdir(CurrentWorkingDir)

	// Mark execution boundaries in logfile : exec stop
	if e := markLogFile("End"); e != nil {
		fmt.Println(e)
		os.Exit(14)
	}
}

func markLogFile(boundary string) error {
	base := filepath.Join(os.Getenv("HOME"), ".local", "state")
	f, err := os.OpenFile(filepath.Join(base, "pgtools.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	ts := time.Now().Format("2006-01-02 15:04:05")
	f.WriteString(fmt.Sprintf("\n---MARK\n%s : %s\n", ts, boundary))
	f.WriteString("---MARK\n")
	return nil
}
