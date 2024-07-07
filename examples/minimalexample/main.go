package main

import (
	_ "github.com/imthatgin/sas/examples/minimalexample/examplemigrations"
	"github.com/imthatgin/sas/pkg/endpoints"
	"github.com/imthatgin/sas/pkg/migration"
	"github.com/imthatgin/sas/pkg/sas"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	friendlyHeader = "⇨ ${time_rfc3339}\t${level}\t"
	requestHeader  = "⇨ ${time_rfc3339}\tHTTP\t${method} ${uri} -> RESP ${status} (took ${latency_human}) (▼${bytes_in}B  ▲${bytes_out}B)\n"
)

type Entry struct {
	sas.DefaultModel

	Content string
}

func init() {
	migration.MigrationsTableName = "__example_migrations_meta"
}

func main() {
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: requestHeader,
	}))

	e.HideBanner = true

	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetHeader("${time_rfc3339} ${level}")
	}
	log.SetHeader(friendlyHeader)

	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		log.Fatal("In-memory database could not be created: ", err)
	}

	models := addModels(db)
	_ = sas.New(e, db, models)

	log.Fatal(e.Start(":8082"))
}

func setupAuth(db *gorm.DB) {

}

func addModels(db *gorm.DB) []sas.Provider {
	entryPolicy := sas.NewPolicy[Entry](endpoints.AllEndpoints)
	entryPolicy.
		CanListAll(func(c echo.Context) bool {
			return true
		}).
		CanListById(func(c echo.Context, entity Entry) bool {
			return entity.ID%2 == 0
		}).
		CanWriteById(func(c echo.Context, entity Entry) bool {
			return true
		})

	entries := sas.FromModel[Entry]("entries", db, entryPolicy)

	entries.WriteBindType(struct {
		Content string
	}{})

	return []sas.Provider{
		&entries,
	}
}
