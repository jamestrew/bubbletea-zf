package main

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/textinput"
	tea "github.com/charmbracelet/bubbletea/v2"
	fuzzy "github.com/sahilm/fuzzy"
)

const (
	listHeight = 24
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Quit   key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "ctrl+p"),
			key.WithHelp("↑/ctrl+p", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "ctrl+n"),
			key.WithHelp("↓/ctrl+n", "move down"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select file"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "esc"),
			key.WithHelp("ctrl+c/esc", "quit"),
		),
	}
}

type model struct {
	textInput     textinput.Model
	keys          keyMap
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

	tiKeyMap := textinput.DefaultKeyMap()
	tiKeyMap.CharacterForward.SetEnabled(false)
	tiKeyMap.CharacterBackward.SetEnabled(false)
	tiKeyMap.WordForward.SetEnabled(false)
	tiKeyMap.WordBackward.SetEnabled(false)
	ti.KeyMap = tiKeyMap

	return model{
		textInput:     ti,
		keys:          newKeyMap(),
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
		slog.Debug("Key pressed", "key", msg.String())

		switch {
		case key.Matches(msg, m.keys.Quit):
			slog.Info("User quit application")
			return m, tea.Quit

		case key.Matches(msg, m.keys.Select):
			if len(m.filteredFiles) > 0 {
				slog.Info("File selected", "file", m.filteredFiles[m.cursor], "cursor", m.cursor)
			}
			return m, tea.Quit

		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
				slog.Debug("Cursor moved up", "cursor", m.cursor)
			}
			return m, nil

		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.filteredFiles)-1 && m.cursor < listHeight-1 {
				m.cursor++
				slog.Debug("Cursor moved down", "cursor", m.cursor)
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		slog.Debug("Window resized", "width", msg.Width, "height", msg.Height)
		m.width = msg.Width
		m.height = msg.Height
	}

	m.textInput, cmd = m.textInput.Update(msg)

	if m.textInput.Value() != "" {
		m.query = m.textInput.Value()
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

	slog.Debug("Filter updated",
		"query", m.query,
		"totalMatches", len(matches),
		"displayed", len(m.filteredFiles),
		"cursor", m.cursor)
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
