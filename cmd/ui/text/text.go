package text

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tomasohCHOM/github-stats/cmd/state"
)

var (
	headerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#A2D2FB")).Bold(true)
	inputStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#ECF2F8")).Bold(true)
	errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FA7970")).Bold(true)
)

type (
	errMsg error
)

type model struct {
	userOptions *state.UserOptions
	textInput   textinput.Model
	err         error
	errMsg      string
	header      string
}

func InitialTextModel(userOptions *state.UserOptions, header, errMsg string) model {
	ti := textinput.New()
	ti.Placeholder = "tomasohCHOM"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		userOptions: userOptions,
		textInput:   ti,
		err:         nil,
		errMsg:      errMsg,
		header:      header,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			m.textInput.Blur()
			m.userOptions.Username = m.textInput.Value()
			return m, tea.Quit
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n%s\n\n%s",
		errorStyle.Render(m.errMsg),
		headerStyle.Render(m.header),
		inputStyle.Render(m.textInput.View()),
	)
}
