package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
)


const (
	listHeight = 10
)

type model struct {
	allFiles     []string
	filteredFiles []string
	query        string
	cursor       int
	width        int
	height       int
}

func initialModel(files []string) model {
	return model{
		allFiles:      files,
		filteredFiles: files[:min(listHeight, len(files))],
		query:         "",
		cursor:        0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			// TODO: handle file selection
			return m, tea.Quit

		case "up", "ctrl+k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "ctrl+j":
			if m.cursor < len(m.filteredFiles)-1 && m.cursor < listHeight-1 {
				m.cursor++
			}

		case "backspace":
			if len(m.query) > 0 {
				m.query = m.query[:len(m.query)-1]
				m.updateFilter()
			}

		default:
			// Handle regular character input
			key := msg.String()
			if len(key) == 1 {
				m.query += key
				m.updateFilter()
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m *model) updateFilter() {
	// TODO(human): Implement fuzzy filtering logic here
	// For now, just do simple substring matching
	m.filteredFiles = []string{}
	query := strings.ToLower(m.query)

	for _, file := range m.allFiles {
		if strings.Contains(strings.ToLower(file), query) {
			m.filteredFiles = append(m.filteredFiles, file)
			if len(m.filteredFiles) >= listHeight {
				break
			}
		}
	}

	// Reset cursor if out of bounds
	if m.cursor >= len(m.filteredFiles) {
		m.cursor = max(0, len(m.filteredFiles)-1)
	}
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("> %s\n", m.query))
	b.WriteString(strings.Repeat("─", 50) + "\n")

	for i := range listHeight {
		if i < len(m.filteredFiles) {
			cursor := " "
			if i == m.cursor {
				cursor = ">"
			}
			b.WriteString(fmt.Sprintf("%s %s\n", cursor, m.filteredFiles[i]))
		} else {
			b.WriteString("  \n")
		}
	}

	b.WriteString("\n")
	b.WriteString("ctrl+c/esc: quit | enter: select | ↑↓: navigate\n")

	return b.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
