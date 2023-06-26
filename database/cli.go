package database

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"github.com/erfansahebi/lamia_auth/config"
	"github.com/erfansahebi/lamia_shared/log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/erfansahebi/lamia_auth/di"
	goMigrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
)

func Migrate(ctx context.Context, configuration *config.Config, steps int) error {
	db, err := sql.Open("postgres", configuration.GetDbUrl("migrate"))
	if err != nil {
		log.WithError(err).Fatalf(ctx, "migrate: failed to open db connection")
		return err
	}

	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		log.WithError(err).Fatalf(ctx, "migrate: failed to get db driver")
		return err
	}

	m, err := goMigrate.NewWithDatabaseInstance(
		"file://"+configuration.Migration.Directory,
		(*configuration).Database.Name, driver)
	if err != nil {
		log.WithError(err).Fatalf(ctx, "migrate: failed to get migrate instances")
		return err
	}

	m.Log = &migrateLogger{IsVerbose: true}

	switch steps {
	case 0:
		err = m.Up()
	case -1:
		err = m.Drop()
	default:
		err = m.Steps(steps)
	}

	switch {
	case err == nil:
		log.Infof(ctx, "Successfully migrated")
	case errors.Is(err, goMigrate.ErrNoChange):
		log.Infof(ctx, "No changes")
	default:
		log.WithError(err).Fatalf(ctx, "failed to migrate")
	}

	return err
}

func nextSeqVersion(matches []string, seqDigits int) (string, error) {
	nextSeq := uint64(1)

	if len(matches) > 0 {
		filename := matches[len(matches)-1]
		matchSeqStr := filepath.Base(filename)
		idx := strings.Index(matchSeqStr, "_")

		if idx < 1 { // Using 1 instead of 0 since there should be at least 1 digit
			return "", fmt.Errorf("malformed migration filename: %s", filename)
		}

		var err error
		matchSeqStr = matchSeqStr[0:idx]
		nextSeq, err = strconv.ParseUint(matchSeqStr, 10, 64)

		if err != nil {
			return "", err
		}

		nextSeq++
	}

	version := fmt.Sprintf("%0[2]*[1]d", nextSeq, seqDigits)

	if len(version) > seqDigits {
		return "", fmt.Errorf("next sequence number %s too large. at most %d digits are allowed", version, seqDigits)
	}

	return version, nil
}

func createFile(filename string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)

	if err != nil {
		return err
	}

	return f.Close()
}

func MakeMigration(ctx context.Context, configuration *config.Config, name string) error {
	diContainer := di.NewDIContainer(ctx, configuration)
	var version string
	var err error
	dir := diContainer.Config().Migration.Directory
	dir = filepath.Clean(dir)

	matches, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return err
	}

	version, err = nextSeqVersion(matches, diContainer.Config().Migration.SeqDigits)
	if err != nil {
		return err
	}

	versionGlob := filepath.Join(dir, version+"_*.sql")
	matches, err = filepath.Glob(versionGlob)

	if err != nil {
		return err
	}

	if len(matches) > 0 {
		return fmt.Errorf("duplicate migration version: %s", version)
	}

	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	for _, direction := range []string{"up", "down"} {
		basename := fmt.Sprintf("%s_%s.%s.sql", version, name, direction)
		filename := filepath.Join(dir, basename)

		if err = createFile(filename); err != nil {
			return err
		}

		absPath, _ := filepath.Abs(filename)
		fmt.Println(absPath)
	}

	return nil
}

type migrateLogger struct {
	IsVerbose bool
}

// Printf or use logger
func (l *migrateLogger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

// Verbose should return true when verbose logging output is wanted
func (l *migrateLogger) Verbose() bool {
	return l.IsVerbose
}
