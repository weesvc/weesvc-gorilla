package cmd

import (
	"fmt"
	"log/slog"
	"sort"

	"github.com/weesvc/weesvc-gorilla/internal/config"

	"github.com/spf13/viper"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/weesvc/weesvc-gorilla/internal/app"
	"github.com/weesvc/weesvc-gorilla/internal/migrations"
)

func newMigrateCommand(config *config.Config) *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Readies the application database",
		RunE: func(cmd *cobra.Command, _ []string) error {
			svc, err := app.New(config)
			if err != nil {
				return err
			}
			defer func() {
				_ = svc.Close()
			}()
			return migrate(cmd, svc)
		},
	}

	migrateCmd.Flags().Int("number", -1, "the migration to run forwards until; if not set, will run all migrations")
	migrateCmd.Flags().Bool("dry-run", false, "print out migrations to be applied without running them")

	migrateCmd.PersistentFlags().StringVar(&config.Dialect, "dialect", "sqlite3", "config file")
	migrateCmd.PersistentFlags().StringVar(&config.DatabaseURI, "database-uri", "", "config file")

	_ = viper.BindPFlag("Dialect", migrateCmd.PersistentFlags().Lookup("dialect"))
	_ = viper.BindPFlag("DatabaseURI", migrateCmd.PersistentFlags().Lookup("database-uri"))

	return migrateCmd
}

func migrate(cmd *cobra.Command, app *app.App) error {
	number, _ := cmd.Flags().GetInt("number")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	if dryRun {
		slog.Info("=== DRY RUN ===")
	}

	runMigrations, err := shouldRunMigrations(app, number)
	if err != nil {
		return err
	}
	if !runMigrations {
		return nil
	}

	for _, migration := range migrations.Migrations {
		if migration.Number > uint(number) {
			break
		}

		logger := slog.With("migration_number", migration.Number)
		logger.Info(fmt.Sprintf("applying migration %q", migration.Name))

		if dryRun {
			continue
		}
		if ok := applyMigration(app, logger, migration); !ok {
			break
		}
	}

	return nil
}

func applyMigration(app *app.App, logger *slog.Logger, migration *migrations.Migration) bool {
	tx := app.Database.Begin()
	if err := migration.Forwards(tx); err != nil {
		logger.With("err", err).Error("unable to apply migration, rolling back")
		if err := tx.Rollback().Error; err != nil {
			logger.With("err", err).Error("unable to rollback...")
		}
		return false
	}

	if err := tx.Commit().Error; err != nil {
		logger.With("err", err).Error("unable to commit transaction")
		return false
	}

	if err := app.Database.Create(migration).Error; err != nil {
		logger.With("err", err).Error("unable to create migration record")
		return false
	}
	return true
}

func shouldRunMigrations(app *app.App, number int) (bool, error) {
	sort.Slice(migrations.Migrations, func(i, j int) bool {
		return migrations.Migrations[i].Number < migrations.Migrations[j].Number
	})

	// Make sure Migration table is there
	if err := app.Database.AutoMigrate(&migrations.Migration{}).Error; err != nil {
		return false, errors.Wrap(err, "unable to automatically migrate migrations table")
	}

	var latest migrations.Migration
	if err := app.Database.Order("number desc").First(&latest).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return false, errors.Wrap(err, "unable to find latest migration")
	}

	noMigrationsApplied := latest.Number == 0
	if (noMigrationsApplied && len(migrations.Migrations) == 0) ||
		(latest.Number >= migrations.Migrations[len(migrations.Migrations)-1].Number) {
		slog.Info("no migrations to apply")
		return false, nil
	}

	if number == -1 {
		number = int(migrations.Migrations[len(migrations.Migrations)-1].Number)
	}
	if uint(number) <= latest.Number && latest.Number > 0 {
		slog.Info("no migrations to apply; number is less than or equal to latest migration")
		return false, nil
	}

	return true, nil
}
