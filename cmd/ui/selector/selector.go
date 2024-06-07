package selector

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
	dimStyle              = lipgloss.NewStyle().Foreground(lipgloss.Color("#C6Cdd5")).Bold(true)
)

type model struct {
	userOptions *program.ProgramState
	options     []string
	cursor      int
	selected    int
	err         error
	header      string
}

func InitialSelectionModel(userOptions *program.ProgramState, header string, options []string) model {
	return model{
		userOptions: userOptions,
		header:      header,
		options:     options,
		selected:    0,
		cursor:      0,
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
		case "q", "ctrl+c", "esc":
			m.userOptions.ExitState = true
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
			m.selected = m.cursor

		case "enter":
			m.userOptions.SelectedContinueOption, m.err = m.handleSelection()
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
		prefix := "( )"
		if i == m.selected {
			prefix = selectedCheckboxStyle.Render("(•)")
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
	helpOptions := "(Press [space] to select, enter to continue. Press q, esc, or ctrl-c to quit)"
	s.WriteString(fmt.Sprintf("\n%s\n", helpOptions))

	return s.String()
}

func (m *model) handleSelection() (string, error) {
	for i := range m.options {
		if i == m.selected {
			return m.options[i], nil
		}
	}
	return "", fmt.Errorf("no options selected")
}
