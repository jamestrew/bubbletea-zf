package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/v2/textinput"
	tea "github.com/charmbracelet/bubbletea/v2"
	fuzzy "github.com/sahilm/fuzzy"
)

const (
	listHeight = 24
)

type model struct {
	textInput     textinput.Model
	allFiles      []string
	filteredFiles []string
	query         string
	cursor        int
	width         int
	height        int
}

func initialModel(files []string) model {
	ti := textinput.New()
	ti.Placeholder = "Search files..."
	ti.Focus()

	return model{
		textInput:     ti,
		allFiles:      files,
		filteredFiles: files[:min(listHeight, len(files))],
		cursor:        0,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

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
			return m, nil

		case "down", "ctrl+j":
			if m.cursor < len(m.filteredFiles)-1 && m.cursor < listHeight-1 {
				m.cursor++
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.textInput, cmd = m.textInput.Update(msg)

	if m.textInput.Value() != "" {
		m.updateFilter()
	} else {
		m.filteredFiles = m.allFiles[:min(listHeight, len(m.allFiles))]
		m.cursor = 0
	}

	return m, cmd
}

func (m *model) updateFilter() {
	m.filteredFiles = []string{}
	matches := fuzzy.Find(m.query, m.allFiles)

	resultsCount := min(len(matches), listHeight)
	m.filteredFiles = make([]string, 0, resultsCount)
	for i := range resultsCount {
		m.filteredFiles = append(m.filteredFiles, matches[i].Str)
	}

	if m.cursor >= len(m.filteredFiles) {
		m.cursor = max(0, len(m.filteredFiles)-1)
	}
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString(m.textInput.View() + "\n")
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
