package cmd

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/httpfs"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"

	"github.com/spf13/cobra"
	"github.com/vamika-digital/wms-api-server/config"
)

//go:embed migrations/*
var migrationFiles embed.FS

func createMigrateInstance() (*migrate.Migrate, error) {
	fs := http.FS(migrationFiles)
	httpFs, err := httpfs.New(fs, "migrations")
	if err != nil {
		return nil, fmt.Errorf("error creating migration source: %v", err)
	}

	databaseURL := config.AppConfig.GetDatabaseConnectionString()
	sourceURL := fmt.Sprintf("httpfs://%s", httpFs)

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	return m, nil
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "A collection of commands to manage database migrations",
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Migrate the database to the most recent version available",
	Run:   runUpMigrations,
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Revert the previous database migration",
	Run:   runDownMigrations,
}

func init() {
	migrateCmd.AddCommand(upCmd)
	migrateCmd.AddCommand(downCmd)
	rootCmd.AddCommand(migrateCmd)
}

func runUpMigrations(cmd *cobra.Command, args []string) {
	m, err := createMigrateInstance()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if err := m.Up(); err != nil {
		log.Fatalf("failed to migrate up: %v", err)
	}
}

func runDownMigrations(cmd *cobra.Command, args []string) {
	m, err := createMigrateInstance()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if err := m.Steps(-1); err != nil {
		log.Fatalf("failed to migrate up: %v", err)
	}
}
