package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lazydb",
	Short: "A terminal UI for your database",
	Long: `lazydb is a fast, keyboard-driven terminal UI for Postgres, MySQL and SQLite.
Browse tables, run queries and inspect your database — all from your terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("lazydb starting...")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}