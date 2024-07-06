package sas

import (
	"github.com/imthatgin/sas/pkg/endpoints"
	"github.com/labstack/echo/v4"
)

// Policy represents the permission requirements and endpoints enabled for a given model.
type Policy[T any] struct {
	EnabledEndpoints endpoints.HttpEndpointType

	canListAll    func(c echo.Context) bool
	canListById   func(c echo.Context, entity T) bool
	canWriteById  func(c echo.Context, entity T) bool
	canCreate     func(c echo.Context) bool
	canDeleteById func(c echo.Context, entity T) bool
}

// NewPolicy creates a default policy instance, which will deny all operations by default.
// Use the exposed methods on the Policy type to set the predicates.
func NewPolicy[T any](endpoints endpoints.HttpEndpointType) Policy[T] {
	return Policy[T]{
		EnabledEndpoints: endpoints,

		canListAll: func(c echo.Context) bool {
			return false
		},

		canListById: func(c echo.Context, entity T) bool {
			return false
		},

		canCreate: func(c echo.Context) bool {
			return false
		},

		canDeleteById: func(c echo.Context, entity T) bool {
			return false
		},
	}
}

// CanListAll takes a predicate and determines whether the operation can proceed.
func (p *Policy[T]) CanListAll(predicate func(c echo.Context) bool) *Policy[T] {
	p.canListAll = predicate

	return p
}

// CanListById takes a predicate and determines whether the operation can proceed.
func (p *Policy[T]) CanListById(predicate func(c echo.Context, entity T) bool) *Policy[T] {
	p.canListById = predicate

	return p
}

// CanWriteById takes a predicate and determines whether the operation can proceed.
func (p *Policy[T]) CanWriteById(predicate func(c echo.Context, entity T) bool) *Policy[T] {
	p.canWriteById = predicate

	return p
}

// CanDeleteById takes a predicate and determines whether the operation can proceed.
func (p *Policy[T]) CanDeleteById(predicate func(c echo.Context, entity T) bool) *Policy[T] {
	p.canDeleteById = predicate

	return p
}
