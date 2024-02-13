package cli

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/database"
	_ "github.com/go-sql-driver/mysql" // for mysql dialect
	_ "github.com/lib/pq"              // for postgres dialect
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
)

// GetMigrationCommandDef initialize migration command
func GetMigrationCommandDef(cfg config.AppConfig) cobra.Command {
	migrateCmd := cobra.Command{
		Use:   "migrate [sub command]",
		Short: "To run db migrate",
		Long: `This command is used to run database migration.
	It has up and down sub commands`,
		Args: cobra.MinimumNArgs(1),
	}

	migrateUp := cobra.Command{
		Use:   "up",
		Short: "It will apply migration(s)",
		Long:  `It will run all remaining migration(s)`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Run test db migration

			if cfg.Env == "testing" && cfg.DB.Dialect == "sqlite3" {
				createTestingDBMigration(cfg, "UP")
			}

			switch cfg.DB.Dialect {
			case database.SQLITE3:
				return runSQLiteMigration(cfg, "UP", false, "")
			case database.POSTGRES:
				return runPostgresMigration(cfg, "UP")
			case database.MYSQL:
				return runMySQLMigration(cfg, "UP")
			}
			return nil
		},
	}

	migrateDown := cobra.Command{
		Use:   "down",
		Short: "It will revert migration(s)",
		Long:  `It will run all remaining migration(s)`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Run test db migration
			if cfg.Env == "testing" && cfg.DB.Dialect == "sqlite3" {
				createTestingDBMigration(cfg, "DOWN")
			}

			switch cfg.DB.Dialect {
			case database.SQLITE3:
				return runSQLiteMigration(cfg, "DOWN", false, "")
			case database.POSTGRES:
				return runPostgresMigration(cfg, "DOWN")
			case database.MYSQL:
				return runMySQLMigration(cfg, "DOWN")
			}
			return nil
		},
	}
	migrateCmd.AddCommand(&migrateUp, &migrateDown)
	// Migration commands up, down

	return migrateCmd
}

func createTestingDBMigration(cfg config.AppConfig, migrationType string) {
	if cfg.Env == "local" || cfg.Env == "testing" {
		err := runSQLiteMigration(cfg, migrationType, true, "database/go-test-db.db")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func runSQLiteMigration(cfg config.AppConfig, migrationType string, testDB bool, testDBPath string) error {
	var migrations = migrate.FileMigrationSource{
		Dir: cfg.DB.MigrationDir,
	}
	var dbPath = cfg.DB.SQLiteFilePath
	if testDB {
		dbPath = testDBPath
	}

	db, err := sql.Open(database.SQLITE3, dbPath)
	if err != nil {
		return err
	}

	if migrationType == "UP" {
		_, err = migrate.Exec(db, database.SQLITE3, migrations, migrate.Up)
		if err != nil {
			return err
		}
	} else {
		_, err = migrate.Exec(db, database.SQLITE3, migrations, migrate.Down)
		if err != nil {
			return err
		}
	}

	return nil
}

func runMySQLMigration(cfg config.AppConfig, migrationType string) error {
	migrations := migrate.FileMigrationSource{
		Dir: cfg.DB.MigrationDir,
	}

	db, err := sql.Open(database.MYSQL, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Db, cfg.DB.QueryString))
	if err != nil {
		return err
	}

	if migrationType == "UP" {
		_, err = migrate.Exec(db, database.MYSQL, migrations, migrate.Up)
		if err != nil {
			return err
		}
	} else {
		_, err = migrate.Exec(db, database.MYSQL, migrations, migrate.Down)
		if err != nil {
			return err
		}
	}

	return nil
}

func runPostgresMigration(cfg config.AppConfig, migrationType string) error {
	migrations := migrate.FileMigrationSource{
		Dir: cfg.DB.MigrationDir,
	}

	db, err := sql.Open(database.POSTGRES, fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Db, cfg.DB.QueryString))
	if err != nil {
		return err
	}

	if migrationType == "UP" {
		_, err = migrate.Exec(db, database.POSTGRES, migrations, migrate.Up)
		if err != nil {
			return err
		}
	} else {
		_, err = migrate.Exec(db, database.POSTGRES, migrations, migrate.Down)
		if err != nil {
			return err
		}
	}

	return nil
}
