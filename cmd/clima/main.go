package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/esferadigital/clima/internal/tui"
)

const DEBUG_PATH = "dev/debug.log"

func main() {
	var (
		sink *os.File
		err  error
	)

	debug := flag.Bool("debug", false, "Save logs to file")
	flag.Parse()

	if *debug {
		if err = os.MkdirAll(filepath.Dir(DEBUG_PATH), os.ModePerm); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to ensure debug directory exists: %v\n", err)
			os.Exit(1)
		}
		if sink, err = os.OpenFile(DEBUG_PATH, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open debug log file: %v\n", err)
			os.Exit(1)
		}
		defer sink.Close()
	}

	if _, err = tea.NewProgram(tui.InitialModel(sink), tea.WithAltScreen()).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "TUI program run failed: %v\n", err)
		os.Exit(1)
	}
}

