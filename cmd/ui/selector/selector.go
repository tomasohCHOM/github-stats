// Define selection options
package selector

import (
	"fmt"
	"strings"

	bubbletea "github.com/charmbracelet/bubbletea"
	// "github.com/tomasohCHOM/github-stats/cmd/stats"
)

var choices = []string{"Taro", "Coffee", "Lychee"}

type model struct {
	options  []string
	cursor   int
	selected map[int]bool
	input    string
	result   string
	err      error
	header   string
}

func InitialSelectionModel(header string, options []string) model {
	return model{
		header:   header,
		options:  options,
		cursor:   0,
		selected: make(map[int]bool),
		input:    "",
		result:   "",
		err:      nil,
	}
}

func (m model) Init() bubbletea.Cmd {
	return nil
}

func (m model) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
	switch msg := msg.(type) {
	case bubbletea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, bubbletea.Quit

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}

		case "down", "j":
			m.cursor++
			if m.cursor >= len(choices) {
				m.cursor = 0
			}

		case " ":
			// Send the choice on the channel and exit.
			m.selected[m.cursor] = !m.selected[m.cursor]
		case "enter":
			m.result, m.err = m.handleSelection()
			return m, bubbletea.Quit
		}

	}
	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}
	s.WriteString(m.header)

	for i, choice := range m.options {
		prefix := "[ ]"
		if m.selected[i] {
			prefix = "[X]"
		}

		if i == m.cursor {
			s.WriteString(fmt.Sprintf("> %s %s\n", prefix, choice))
		} else {
			s.WriteString(fmt.Sprintf("  %s %s\n", prefix, choice))
		}
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}

func (m *model) handleSelection() (string, error) {
	var results []string
	for i := range m.selected {
		if m.selected[i] {
			switch m.options[i] {
			case "Pull Requests":
				results = append(results, "Fetching pull requests...") // Add actual PR fetching logic here
			case "Issues":
				results = append(results, "Fetching issues...") // Add actual issues fetching logic here
			}
		}
	}
	if len(results) == 0 {
		return "", fmt.Errorf("no options selected")
	}
	return strings.Join(results, "\n"), nil
}
