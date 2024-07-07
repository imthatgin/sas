package auth

import (
	"github.com/imthatgin/sas/pkg/migration"
	"github.com/imthatgin/sas/pkg/sas"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type OnLoginEvent func(c echo.Context)

type AuthUser struct {
	sas.DefaultModel

	Provider string
}

func NewAuther(db *gorm.DB) *Auther {
	migration.RegisterSystemMigration("add_auth_user", func(db *gorm.DB) error {
		return db.AutoMigrate(AuthUser{})
	}, nil)

	return &Auther{
		db:           db,
		onLoginEvent: []OnLoginEvent{},
	}
}

type Auther struct {
	db *gorm.DB

	onLoginEvent []OnLoginEvent
}

func (a *Auther) OnLogin(handler OnLoginEvent) {
	a.onLoginEvent = append(a.onLoginEvent, handler)
}

func (a *Auther) Register(e *echo.Echo) {
	auth := e.Group("/auth")

	auth.GET("/:provider/login", func(c echo.Context) error {

		return BeginAuthHandler(c)
	})
}
