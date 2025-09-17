package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/esferadigital/clima/tui"
)

func main() {
	debug := flag.Bool("debug", false, "dump logs to file")
	flag.Parse()
	if *debug {
		debugPath := "dev/debug.log"
		if err := os.MkdirAll(filepath.Dir(debugPath), os.ModePerm); err != nil {
			fmt.Printf("Failed to ensure debug directory exists: %v\n", err)
			os.Exit(1)
		}
		if logFile, err := tea.LogToFile("dev/debug.log", "debug"); err != nil {
			fmt.Printf("Failed to set up debug log file: %v\n", err)
			os.Exit(1)
		} else {
			defer logFile.Close()
		}
	}

	if _, err := tea.NewProgram(tui.InitialModel()).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}

