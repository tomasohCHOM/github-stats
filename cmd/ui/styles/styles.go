package styles

import "github.com/charmbracelet/lipgloss"

var (
	LogoStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("#77BDFB")).Bold(true)
	HeaderStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("#A2D2FB")).Bold(true)
	ContrastStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#CEA5FB")).Bold(true)
	StatsStyle            = lipgloss.NewStyle().Foreground(lipgloss.Color("#7CE38B")).Bold(true)
	SelectedCheckboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7CE38B")).Bold(true)
	SelectedTextStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#ECF2F8")).Bold(true)
	BlurStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("#89929B")).Bold(true)
	DimStyle              = lipgloss.NewStyle().Foreground(lipgloss.Color("#C6CDD5")).Bold(true)
)
