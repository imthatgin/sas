package sas

import (
	"github.com/imthatgin/sas/pkg/migration"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type Provider interface {
	Register(e *echo.Echo)
}

type Server struct {
	e  *echo.Echo
	db *gorm.DB

	// Reference database models
	resources []Provider
}

func New(echo *echo.Echo, db *gorm.DB, resources []Provider) *Server {
	echo.HTTPErrorHandler = ManagedModelErrorHandler

	s := &Server{
		e:  echo,
		db: db,

		resources: resources,
	}

	for _, resource := range s.resources {
		resource.Register(echo)
	}

	err := migration.RunMigrations(s.db)
	if err != nil {
		log.Error("Could not run migrations: ", err)
	}

	return s
}
