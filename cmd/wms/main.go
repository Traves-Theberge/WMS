package main

import (
	"flag"
	"fmt"
	"os"

	"wms/internal/config"
	"wms/internal/ui/models"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Parse command line flags
	location := flag.String("location", "", "Location to get weather for")
	units := flag.String("units", "metric", "Units (metric, imperial)")
	timeFormat := flag.String("time", "24", "Time format (12, 24)")
	compact := flag.Bool("compact", false, "Compact display mode")
	refresh := flag.Int("refresh", 5, "Refresh interval in minutes")
	help := flag.Bool("help", false, "Show help")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	// Load configuration
	cfg := config.ReadConfig()

	// Override config with command line flags
	if *location != "" {
		cfg.Location = *location
	}
	if *units != "" {
		cfg.Units = *units
	}
	if *timeFormat != "" {
		cfg.TimeFormat = *timeFormat
	}
	if *compact {
		cfg.Compact = *compact
	}
	if *refresh > 0 {
		cfg.RefreshInterval = *refresh
	}

	// Initialize the model with configuration
	m := models.InitialModelWithConfig(cfg)

	// Create the Bubble Tea program
	p := tea.NewProgram(m, tea.WithAltScreen())

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}