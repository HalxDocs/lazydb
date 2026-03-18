package cmd

import (
	"fmt"
	"os"

	"github.com/HalxDocs/lazydb/internal/app"
	"github.com/HalxDocs/lazydb/internal/config"
	"github.com/spf13/cobra"
)

var (
	driver  string
	dsn     string
	connName string
	saveName string
)

var rootCmd = &cobra.Command{
	Use:   "lazydb",
	Short: "A terminal UI for your database",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		// load a saved connection by name
		if connName != "" {
			conn, err := cfg.Find(connName)
			if err != nil {
				return err
			}
			return app.Run(conn.Driver, conn.DSN)
		}

		if dsn == "" {
			// list saved connections if any exist
			if len(cfg.Connections) > 0 {
				fmt.Println("saved connections:")
				for _, c := range cfg.Connections {
					fmt.Printf("  %-20s %s\n", c.Name, c.Driver)
				}
				fmt.Println("\nuse: lazydb --conn <name>")
				return nil
			}
			return fmt.Errorf(
				"--dsn is required\n\nexamples:\n" +
				"  lazydb --driver postgres --dsn \"postgres://user:pass@localhost:5432/mydb\"\n" +
				"  lazydb --driver sqlite  --dsn ./mydb.sqlite\n\n" +
				"save a connection:\n" +
				"  lazydb --driver postgres --dsn \"...\" --save myapp",
			)
		}

		// save connection if requested
		if saveName != "" {
			cfg.Add(config.Connection{
				Name:   saveName,
				Driver: driver,
				DSN:    dsn,
			})
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("saving connection: %w", err)
			}
			fmt.Printf("connection %q saved\n", saveName)
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
	rootCmd.Flags().StringVar(&connName, "conn", "", "name of a saved connection")
	rootCmd.Flags().StringVar(&saveName, "save", "", "save this connection with a name")
}