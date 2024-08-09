package migration

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

var (
	RegisteredMigrations []Migration
	MigrationsTableName  = "__migrations_meta"
)

type MigratorFunc func(db *gorm.DB) error

type Migration struct {
	Name     string
	FullPath string

	NoChecksum bool

	Up   MigratorFunc
	Down MigratorFunc
}

func Register(up MigratorFunc, down MigratorFunc) {
	// Convert the file name into a migration name.
	_, filePath, _, _ := runtime.Caller(1) // skip 1 for call site stack level for consumer
	_, fileName := path.Split(filePath)
	migrationName := strings.Split(fileName, ".")[0]

	RegisteredMigrations = append(RegisteredMigrations, Migration{
		Name:     migrationName,
		FullPath: filePath,
		Up:       up,
		Down:     down,
	})
}

func RegisterSystemMigration(name string, up MigratorFunc, down MigratorFunc) {
	RegisteredMigrations = append(RegisteredMigrations, Migration{
		Name:       name,
		NoChecksum: true,
		Up:         up,
		Down:       down,
	})
}

type MigrationsMeta struct {
	Name     string
	Checksum string

	Timestamp time.Time
}

func RunMigrations(db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		log.Infof("Ensuring meta table exists")
		metaErr := tx.Table(MigrationsTableName).AutoMigrate(MigrationsMeta{})
		if metaErr != nil {
			return metaErr
		}

		log.Infof("Will migrate over %d migrations", len(RegisteredMigrations))

		for i, migration := range RegisteredMigrations {
			log.Infof("(%d) Migrating %s", i, migration.Name)

			var sum string

			// Check if we have a matching migration with a correct checksum, and skip if that is the case.
			var existing []MigrationsMeta
			tx.Table(MigrationsTableName).Limit(1).Find(&existing, "name = ?", migration.Name)
			if len(existing) > 0 {
				// Allows system migrations to ignore checksums
				if !migration.NoChecksum {
					sum, _ = getMigrationChecksum(migration)
					if sum == "" {
						return fmt.Errorf("checksum was empty for migration %s", migration.Name)
					}
					if existing[0].Checksum != sum {
						return fmt.Errorf("checksum mismatch for migration %s: %s != %s", migration.Name, sum, existing[0].Checksum)
					}
				} else {
					sum = time.Now().UTC().String()
				}

				log.Infof(">\t DONE skip %s", migration.Name)
				continue
			} else {
				sum, _ = getMigrationChecksum(migration)
			}

			err := migration.Up(tx)
			if err != nil {
				log.Errorf(">\t FAIL migration %s: %s", migration.Name, err)
				return err
			}

			// Insert the migration into the meta table
			timeStamp := time.Now().UTC()
			tx.Table(MigrationsTableName).Create(MigrationsMeta{
				Name:      migration.Name,
				Checksum:  sum,
				Timestamp: timeStamp,
			})
			log.Infof("Added migration meta for %s with checksum %s", migration.Name, sum)

			log.Infof(">\t DONE migration %s", migration.Name)
		}
		return nil
	})

	return err
}

// getMigrationChecksum will get the MD5 hash of the migration specified.
// This is used to avoid applying migrations incorrectly.
func getMigrationChecksum(migration Migration) (string, error) {
	file, err := os.Open(migration.FullPath)
	if err != nil {
		return "", err
	}

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	err = file.Close()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
