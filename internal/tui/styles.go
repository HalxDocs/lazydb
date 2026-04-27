package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Palette
	colorPrimary    = lipgloss.Color("#7EB8F7")
	colorAccent     = lipgloss.Color("#C792EA")
	colorSecondary  = lipgloss.Color("#4A5880")
	colorBackground = lipgloss.Color("#0B0E15")
	colorSurface    = lipgloss.Color("#0E1220")
	colorSurfaceAlt = lipgloss.Color("#11162A")
	colorBorder     = lipgloss.Color("#1A2035")
	colorBorderHi   = lipgloss.Color("#2A3550")
	colorSuccess    = lipgloss.Color("#3ECF8E")
	colorWarning    = lipgloss.Color("#F0A429")
	colorError      = lipgloss.Color("#E05260")
	colorMuted      = lipgloss.Color("#3A4060")
	colorText       = lipgloss.Color("#C8D0EA")
	colorTextDim    = lipgloss.Color("#8892B0")
	colorNumeric    = lipgloss.Color("#F0A429")
	colorString     = lipgloss.Color("#C8D0EA")
	colorNull       = lipgloss.Color("#5C5F7A")
	colorBool       = lipgloss.Color("#3ECF8E")
	colorPK         = lipgloss.Color("#F0A429")

	// Sidebar
	SidebarStyle = lipgloss.NewStyle().
			Background(colorSurface).
			BorderRight(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorBorder).
			PaddingLeft(1).
			PaddingRight(1)

	SidebarTitleStyle = lipgloss.NewStyle().
				Foreground(colorAccent).
				Background(colorSurface).
				Bold(true).
				PaddingLeft(1)

	SidebarSubtitleStyle = lipgloss.NewStyle().
				Foreground(colorMuted).
				Background(colorSurface).
				PaddingLeft(1)

	SidebarItemStyle = lipgloss.NewStyle().
				Foreground(colorTextDim).
				Background(colorSurface).
				PaddingLeft(1)

	SidebarItemActiveStyle = lipgloss.NewStyle().
				Foreground(colorPrimary).
				Background(lipgloss.Color("#1A2540")).
				Bold(true).
				PaddingLeft(1)

	SidebarCountStyle = lipgloss.NewStyle().
				Foreground(colorMuted).
				Background(colorSurface)

	SidebarCountActiveStyle = lipgloss.NewStyle().
				Foreground(colorPrimary).
				Background(lipgloss.Color("#1A2540"))

	SidebarSearchStyle = lipgloss.NewStyle().
				Foreground(colorPrimary).
				Background(colorSurfaceAlt).
				PaddingLeft(1).
				PaddingRight(1)

	// Toolbar / breadcrumb
	ToolbarStyle = lipgloss.NewStyle().
			Background(colorSurface).
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorBorder).
			PaddingLeft(1)

	BreadcrumbStyle = lipgloss.NewStyle().
			Foreground(colorTextDim).
			Background(colorSurface).
			PaddingLeft(1).
			PaddingRight(1)

	BreadcrumbActiveStyle = lipgloss.NewStyle().
				Foreground(colorPrimary).
				Background(colorSurface).
				Bold(true)

	BreadcrumbSepStyle = lipgloss.NewStyle().
				Foreground(colorMuted).
				Background(colorSurface)

	KbdStyle = lipgloss.NewStyle().
			Foreground(colorText).
			Background(colorBorder).
			Padding(0, 1)

	// Table
	TableHeaderStyle = lipgloss.NewStyle().
				Foreground(colorTextDim).
				Background(colorSurfaceAlt).
				Bold(true)

	TableHeaderActiveStyle = lipgloss.NewStyle().
				Foreground(colorAccent).
				Background(lipgloss.Color("#1A2540")).
				Bold(true)

	TableCellStyle = lipgloss.NewStyle().
			Foreground(colorTextDim)

	TableCellActiveStyle = lipgloss.NewStyle().
				Foreground(colorText).
				Background(lipgloss.Color("#141E30"))

	TableCellFocusedStyle = lipgloss.NewStyle().
				Foreground(colorBackground).
				Background(colorPrimary).
				Bold(true)

	TableCellNumericStyle = lipgloss.NewStyle().
				Foreground(colorNumeric)

	TableCellNullStyle = lipgloss.NewStyle().
				Foreground(colorNull).
				Italic(true)

	TableCellBoolStyle = lipgloss.NewStyle().
				Foreground(colorBool)

	// Query bar
	QueryBarStyle = lipgloss.NewStyle().
			Background(colorSurface).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorBorderHi).
			PaddingLeft(1).
			PaddingRight(1)

	QueryPromptStyle = lipgloss.NewStyle().
				Foreground(colorAccent).
				Bold(true)

	QueryInputStyle = lipgloss.NewStyle().
				Foreground(colorPrimary)

	FilterBarStyle = lipgloss.NewStyle().
			Background(colorSurface).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorWarning).
			PaddingLeft(1).
			PaddingRight(1)

	FilterPromptStyle = lipgloss.NewStyle().
				Foreground(colorWarning).
				Bold(true)

	// Status bar
	StatusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#0A2040")).
			Foreground(colorText)

	StatusBarTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7EB8F7")).
				Background(lipgloss.Color("#0A2040")).
				Bold(false)

	StatusBarModeStyle = lipgloss.NewStyle().
				Foreground(colorBackground).
				Background(colorPrimary).
				Bold(true).
				Padding(0, 1)

	StatusBarModeQueryStyle = lipgloss.NewStyle().
				Foreground(colorBackground).
				Background(colorAccent).
				Bold(true).
				Padding(0, 1)

	StatusBarModeFilterStyle = lipgloss.NewStyle().
					Foreground(colorBackground).
					Background(colorWarning).
					Bold(true).
					Padding(0, 1)

	StatusBarHintStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#5471A0")).
				Background(lipgloss.Color("#0A2040"))

	// Error / toast
	ErrorStyle = lipgloss.NewStyle().
			Foreground(colorError).
			PaddingLeft(1)

	ToastStyle = lipgloss.NewStyle().
			Foreground(colorBackground).
			Background(colorSuccess).
			Bold(true).
			Padding(0, 1)

	ToastErrorStyle = lipgloss.NewStyle().
			Foreground(colorBackground).
			Background(colorError).
			Bold(true).
			Padding(0, 1)

	// Modal
	ModalStyle = lipgloss.NewStyle().
			Background(colorSurfaceAlt).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary).
			Padding(1, 2)

	ModalTitleStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true).
			Background(colorSurfaceAlt)

	ModalKeyStyle = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Background(colorSurfaceAlt).
			Bold(true)

	ModalValueStyle = lipgloss.NewStyle().
			Foreground(colorText).
			Background(colorSurfaceAlt)

	ModalDimStyle = lipgloss.NewStyle().
			Foreground(colorTextDim).
			Background(colorSurfaceAlt)

	// Schema
	SchemaTypeStyle = lipgloss.NewStyle().
			Foreground(colorAccent)

	SchemaPKStyle = lipgloss.NewStyle().
			Foreground(colorPK).
			Bold(true)

	SchemaNullableStyle = lipgloss.NewStyle().
				Foreground(colorMuted).
				Italic(true)

	// Misc
	TableDividerStyle = lipgloss.NewStyle().
				Foreground(colorBorder)

	RowCursorStyle = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true)

	RowCountStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			PaddingLeft(2)

	SortAscIndicator  = "↑"
	SortDescIndicator = "↓"
)
