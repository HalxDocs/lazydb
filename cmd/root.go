package cmd

import (
	"fmt"
	"os"

	"github.com/HalxDocs/lazydb/internal/app"
	"github.com/spf13/cobra"
)

var (
	driver string
	dsn    string
)

var rootCmd = &cobra.Command{
	Use:   "lazydb",
	Short: "A terminal UI for your database",
	Long: `lazydb is a fast, keyboard-driven terminal UI for Postgres, MySQL and SQLite.
Browse tables, run queries and inspect your database — all from your terminal.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if dsn == "" {
			return fmt.Errorf("--dsn is required\n\nexamples:\n  lazydb --driver postgres --dsn \"postgres://user:pass@localhost:5432/mydb\"\n  lazydb --driver sqlite  --dsn ./mydb.sqlite")
		}
		return app.Run(driver, dsn)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&driver, "driver", "postgres", "database driver: postgres, mysql, sqlite")
	rootCmd.Flags().StringVar(&dsn, "dsn", "", "database connection string")
}