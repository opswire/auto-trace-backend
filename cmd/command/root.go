package main

import (
	"car-sell-buy-system/config"
	"fmt"
	"github.com/fatih/color"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"os"
)

var (
	cfg *config.Config

	successColor = color.New(color.FgGreen).SprintFunc()
	errorColor   = color.New(color.FgRed).SprintFunc()
	warnColor    = color.New(color.FgYellow).SprintFunc()
)

var rootCmd = &cobra.Command{
	Use: "Auto-trace CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Welcome to %s! Use --help for usage.", cfg.App.Name)
	},
}

var migrateUp = &cobra.Command{
	Use:   "migrate-up",
	Short: "Запуск миграций",
	Long:  `Запускает все миграции`,
	//Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		m, err := migrate.New(
			"file://./migrations",
			fmt.Sprintf("%s?sslmode=disable", cfg.Pg.ExposeURL),
		)
		if err != nil {
			fmt.Printf("%s: Connect to databse is failed!: %s\n", errorColor("ERROR"), err)
			return
		}

		if err := m.Up(); err != nil {
			fmt.Printf("%s: Migrations up is failed!: %s\n", errorColor("ERROR"), err)
			return
		}

		fmt.Printf("%s: Migrations up completed!\n", successColor("SUCCESS"))
	},
}

var migrateDown = &cobra.Command{
	Use:   "migrate-down",
	Short: "Откат миграций",
	Long:  `Откатывает все миграции`,
	//Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		m, err := migrate.New(
			"file://./migrations",
			fmt.Sprintf("%s?sslmode=disable", cfg.Pg.ExposeURL),
		)
		if err != nil {
			fmt.Printf("%s: Connect to databse is failed!: %s\n", errorColor("ERROR"), err)
			return
		}

		if err := m.Down(); err != nil {
			fmt.Printf("%s: Migrations down is failed!: %s\n", errorColor("ERROR"), err)
			return
		}

		fmt.Printf("%s: Migrations down completed!\n", successColor("SUCCESS"))
	},
}

func init() {
	cfg = config.NewConfig()

	rootCmd.AddCommand(migrateUp)
	rootCmd.AddCommand(migrateDown)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
