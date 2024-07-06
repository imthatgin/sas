package sas

import (
	"errors"
	patch "github.com/geraldo-labs/merge-struct"
	"github.com/imthatgin/sas/pkg/endpoints"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"net/http"
	"reflect"
	"strconv"
)

type ModelResource[T any] struct {
	Name string

	db *gorm.DB

	Policy  Policy[T]
	Queries Queries[T]

	// Binding
	createBindType any
	writeBindType  any

	createTransformer func(c echo.Context) (*T, error)

	middlewares []echo.MiddlewareFunc
	onRegister  func(e *echo.Echo)
}

func FromModel[T any](name string, db *gorm.DB, policy Policy[T]) ModelResource[T] {
	mr := ModelResource[T]{
		Name: name,

		db:     db,
		Policy: policy,

		// Default queries
		Queries: NewQueries[T](),
	}

	return mr
}

// Register is called automatically by SAS, and will add the configured endpoint behaviours to echo.
func (mr *ModelResource[T]) Register(e *echo.Echo) {
	group := e.Group(mr.Name)

	if endpoints.Has(mr.Policy.EnabledEndpoints, endpoints.GET) {
		group.GET("", mr.getAll, mr.middlewares...)
	}

	if endpoints.Has(mr.Policy.EnabledEndpoints, endpoints.GET) {
		group.GET("/:id", mr.getById, mr.middlewares...)
	}

	if endpoints.Has(mr.Policy.EnabledEndpoints, endpoints.PUT) {
		group.PUT("/:id", mr.writeById, mr.middlewares...)
	}

	if endpoints.Has(mr.Policy.EnabledEndpoints, endpoints.POST) {
		group.POST("", mr.create, mr.middlewares...)
	}

	if endpoints.Has(mr.Policy.EnabledEndpoints, endpoints.DELETE) {
		group.DELETE("/:id", mr.deleteById, mr.middlewares...)
	}

	if mr.onRegister != nil {
		mr.onRegister(e)
	}
}

func (mr *ModelResource[T]) getAll(c echo.Context) error {
	if !mr.Policy.canListAll(c) {
		return ErrorResourceNoAccess
	}

	result, err := mr.Queries.listAllQuery(c, mr.db)
	if err != nil {
		return errors.Join(ErrorDatabaseIssue, err)
	}

	return c.JSON(http.StatusOK, result)
}

func (mr *ModelResource[T]) getById(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return errors.Join(ErrorResourceInvalidID, err)
	}

	result, err := mr.Queries.listByIdQuery(c, mr.db, uint(id))
	if err != nil {
		return errors.Join(ErrorDatabaseIssue, err)
	}

	if !mr.Policy.canListById(c, *result) {
		return ErrorResourceNoAccess
	}

	return c.JSON(http.StatusOK, result)
}

func (mr *ModelResource[T]) writeById(c echo.Context) error {
	// Check that we have a bind type set up already. If not, we must fail the call.
	if mr.writeBindType == nil {
		return ErrorFatalSetupNoBindType
	}

	// Try to instantiate the "DTO" type, and bind to it.
	boundType := reflect.TypeOf(mr.writeBindType)
	boundPtr := reflect.New(boundType)
	bound := boundPtr.Interface()
	if err := c.Bind(bound); err != nil {
		return errors.Join(ErrorResourceInvalidData, err)
	}

	// Parse the ID parameter, or fail.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return errors.Join(ErrorResourceInvalidID, err)
	}

	result, err := mr.Queries.listByIdQuery(c, mr.db, uint(id))
	if err != nil {
		return errors.Join(ErrorDatabaseIssue, err)
	}

	if !mr.Policy.canListById(c, *result) {
		return ErrorResourceNoAccess
	}

	err = mr.Queries.writeByIdQuery(c, mr.db, result, bound)
	if err != nil {
		return errors.Join(ErrorDatabaseIssue, err)
	}

	return c.NoContent(http.StatusOK)
}

func (mr *ModelResource[T]) create(c echo.Context) error {
	if !mr.Policy.canCreate(c) {
		return ErrorResourceNoAccess
	}

	// Patch data onto the structure.
	var model T
	if mr.createTransformer != nil {
		transformedModel, err := mr.createTransformer(c)
		if err != nil {
			return errors.Join(ErrorResourceInvalidData, err)
		}

		if transformedModel != nil {
			model = *transformedModel
		}
	} else {
		if mr.createBindType == nil {
			return ErrorFatalSetupNoBindType
		}

		// Try to instantiate the "DTO" type, and bind to it.
		boundType := reflect.TypeOf(mr.createBindType)
		boundPtr := reflect.New(boundType)
		bound := boundPtr.Interface()
		if err := c.Bind(bound); err != nil {
			log.Error("Binding failed: ", err)
			return err
		}

		_, err := patch.Struct(&model, bound)
		if err != nil {
			log.Error("Patching failed: ", err)
			return err
		}
	}

	tx := mr.db.Create(&model)
	if tx.Error != nil {
		return errors.Join(ErrorDatabaseIssue, tx.Error)
	}

	return c.NoContent(http.StatusOK)
}

func (mr *ModelResource[T]) deleteById(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return errors.Join(ErrorResourceInvalidID, err)
	}

	var result T
	tx := mr.db.First(&result, "id = ?", id)
	if tx.Error != nil {
		return errors.Join(ErrorDatabaseIssue, err)
	}

	if !mr.Policy.canDeleteById(c, result) {
		return ErrorResourceNoAccess
	}

	err = mr.Queries.deleteByIdQuery(c, mr.db, result)
	if err != nil {
		return errors.Join(ErrorDatabaseIssue, err)
	}

	return c.NoContent(http.StatusOK)
}

func (mr *ModelResource[T]) CreateBindType(bt any) {
	mr.createBindType = bt
}

func (mr *ModelResource[T]) WriteBindType(bt any) {
	mr.writeBindType = bt
}

func (mr *ModelResource[T]) OnRegister(handler func(e *echo.Echo)) {
	mr.onRegister = handler
}

func (mr *ModelResource[T]) Middlewares(middlewares ...echo.MiddlewareFunc) {
	mr.middlewares = middlewares
}
