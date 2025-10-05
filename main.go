package main

import (
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/jamestrew/bubbletea-zf/files"
)

func main() {
	logFile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()

	slog.SetDefault(slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	slog.Info("Application started")

	pathsList, err := files.GetWithChannel(".")
	if err != nil {
		fmt.Println("Error getting files:", err)
		return
	}

	slog.Info("Files loaded", "count", len(pathsList))

	p := tea.NewProgram(initialModel(pathsList))
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
