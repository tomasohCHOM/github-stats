// Define selection options
package multiselector

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tomasohCHOM/github-stats/cmd/program"
)

var (
	headerStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("#A2D2FB")).Bold(true)
	selectedCheckboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7CE38B")).Bold(true)
	selectedTextStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#ECF2F8")).Bold(true)
	blurStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("#89929B")).Bold(true)
)

type model struct {
	userOptions *program.ProgramState
	options     []string
	cursor      int
	selected    map[int]bool
	input       string
	err         error
	header      string
}

func InitialMultiSelectModel(userOptions *program.ProgramState, header string, options []string) model {
	return model{
		userOptions: userOptions,
		header:      header,
		options:     options,
		cursor:      0,
		selected:    make(map[int]bool),
		input:       "",
		err:         nil,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.options) - 1
			}

		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.options) {
				m.cursor = 0
			}

		case " ":
			// Send the choice on the channel and exit.
			m.selected[m.cursor] = !m.selected[m.cursor]
		case "enter":
			m.userOptions.SelectedStats, m.err = m.handleSelection()
			return m, tea.Quit
		}

	}
	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}
	header := fmt.Sprintf("%s\n", headerStyle.Render(m.header))
	s.WriteString(header)

	for i, choice := range m.options {
		prefix := "[ ]"
		if m.selected[i] {
			prefix = selectedCheckboxStyle.Render("[✓]")
			choice = selectedTextStyle.Render(choice)
		}

		line := fmt.Sprintf("%s %s", prefix, choice)
		if i == m.cursor {
			s.WriteString(fmt.Sprintf("> %s\n", line))
		} else {
			line = blurStyle.Render(line)
			s.WriteString(fmt.Sprintf(" %s\n", line))
		}
	}
	s.WriteString("\n(press q, esc, or ctrl-c to quit)\n")

	return s.String()
}

func (m *model) handleSelection() ([]string, error) {
	var results []string
	for i := range m.selected {
		if m.selected[i] {
			results = append(results, m.options[i])
		}
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no options selected")
	}
	return results, nil
}
