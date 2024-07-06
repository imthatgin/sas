package examplemigrations

import (
	"fmt"
	"github.com/imthatgin/sas/pkg/migration"
	"github.com/imthatgin/sas/pkg/sas"
	"gorm.io/gorm"
)

type entry001 struct {
	sas.DefaultModel

	Content string
}

func init() {
	migration.Register(func(db *gorm.DB) error {
		err := db.Table("entries").AutoMigrate(entry001{})

		for i := 0; i < 10; i++ {
			tx := db.Table("entries").Create(&entry001{
				Content: fmt.Sprintf("Test %d", i),
			})
			if tx.Error != nil {
				return err
			}
		}

		return err
	}, nil)
}
