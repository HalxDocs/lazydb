package app

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/HalxDocs/lazydb/internal/db"
	"github.com/HalxDocs/lazydb/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		fmt.Print("\033[H\033[2J")
	}
}

func Run(driver, dsn string) error {
	database := db.New(db.Config{
		Driver: driver,
		DSN:    dsn,
	})

	if database == nil {
		return fmt.Errorf("unsupported driver: %s (use postgres, mysql or sqlite)", driver)
	}

	if err := database.Connect(); err != nil {
		return fmt.Errorf("could not connect: %w", err)
	}
	defer database.Close()

	clearScreen()

	model := tui.NewModel(database)
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return nil
}