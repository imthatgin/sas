package examplemigrations

import (
	"github.com/imthatgin/sas/pkg/migration"
	"gorm.io/gorm"
)

type entry002 struct {
	Published bool
}

func init() {
	migration.Register(func(db *gorm.DB) error {
		tx := db.Exec("ALTER TABLE entries ADD published TINYINT")
		return tx.Error
	}, nil)
}
