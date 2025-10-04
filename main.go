package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/jamestrew/bubbletea-zf/files"
)

func main() {
	pathsList, err := files.GetWithChannel(".")
	if err != nil {
		fmt.Println("Error getting files:", err)
		return
	}

	p := tea.NewProgram(initialModel(pathsList))
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
