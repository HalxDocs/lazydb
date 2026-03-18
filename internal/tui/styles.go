package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	colorPrimary    = lipgloss.Color("#7EB8F7")
	colorSecondary  = lipgloss.Color("#4A5880")
	colorBackground = lipgloss.Color("#0B0E15")
	colorSurface    = lipgloss.Color("#0E1220")
	colorBorder     = lipgloss.Color("#1A2035")
	colorSuccess    = lipgloss.Color("#3ECF8E")
	colorWarning    = lipgloss.Color("#F0A429")
	colorError      = lipgloss.Color("#E05260")
	colorMuted      = lipgloss.Color("#3A4060")
	colorText       = lipgloss.Color("#C8D0EA")
	colorTextDim    = lipgloss.Color("#8892B0")

	// Sidebar
	SidebarStyle = lipgloss.NewStyle().
			Width(22).
			Background(lipgloss.Color("#0E1220")).
			BorderRight(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorBorder)

	SidebarTitleStyle = lipgloss.NewStyle().
				Foreground(colorMuted).
				Background(lipgloss.Color("#0E1220")).
				PaddingLeft(1).
				Bold(false)

	SidebarItemStyle = lipgloss.NewStyle().
				Foreground(colorSecondary).
				Background(lipgloss.Color("#0E1220")).
				PaddingLeft(2)

	SidebarItemActiveStyle = lipgloss.NewStyle().
				Foreground(colorPrimary).
				Background(lipgloss.Color("#1A2540")).
				PaddingLeft(2)

	// Toolbar
	ToolbarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#0E1220")).
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorBorder).
			PaddingLeft(1)

	ToolbarItemStyle = lipgloss.NewStyle().
				Foreground(colorSecondary)

	ToolbarItemActiveStyle = lipgloss.NewStyle().
				Foreground(colorPrimary)

	KbdStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Background(lipgloss.Color("#1A2035")).
			PaddingLeft(1).
			PaddingRight(1)

	// Table
	TableHeaderStyle = lipgloss.NewStyle().
				Foreground(colorSecondary).
				Background(lipgloss.Color("#111826")).
				Bold(false)

	TableCellStyle = lipgloss.NewStyle().
			Foreground(colorTextDim)

	TableCellActiveStyle = lipgloss.NewStyle().
				Foreground(colorText).
				Background(lipgloss.Color("#141E30"))

	// Query bar
	QueryBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#0E1220")).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorBorder).
			PaddingLeft(1)

	QueryPromptStyle = lipgloss.NewStyle().
				Foreground(colorMuted)

	QueryInputStyle = lipgloss.NewStyle().
				Foreground(colorPrimary)

	// Status bar
	StatusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#0A2040")).
			PaddingLeft(1).
			PaddingRight(1)

	StatusBarTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#3A7ABF")).
				Background(lipgloss.Color("#0A2040"))

	// Error
	ErrorStyle = lipgloss.NewStyle().
			Foreground(colorError).
			PaddingLeft(1)
)