package sas

import (
	"errors"
	patch "github.com/geraldo-labs/merge-struct"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Queries[T any] struct {
	listByIdQuery   func(c echo.Context, q *gorm.DB, id uint) (*T, error)
	listAllQuery    func(c echo.Context, q *gorm.DB) ([]T, error)
	writeByIdQuery  func(c echo.Context, q *gorm.DB, entity *T, new any) error
	deleteByIdQuery func(c echo.Context, q *gorm.DB, entity T) error
}

// NewQueries returns a new instance of the query functions used by default.
// The Queries struct has methods to override each query type.
func NewQueries[T any]() Queries[T] {
	return Queries[T]{
		// TODO: Add pagination options
		listAllQuery: func(c echo.Context, q *gorm.DB) ([]T, error) {
			var result []T
			tx := q.Find(&result)

			if tx.Error != nil {
				return nil, ErrorResourceNotFound
			}

			return result, nil
		},

		listByIdQuery: func(c echo.Context, q *gorm.DB, id uint) (*T, error) {
			var result T
			tx := q.First(&result, "id = ?", id)

			if tx.Error != nil {
				if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
					return nil, ErrorResourceNotFound
				}

				return nil, tx.Error
			}

			return &result, nil
		},

		writeByIdQuery: func(c echo.Context, q *gorm.DB, entity *T, new any) error {
			_, err := patch.Struct(entity, new)
			if err != nil {
				return err
			}

			tx := q.Save(entity)
			if tx.Error != nil {
				return tx.Error
			}

			return nil
		},

		deleteByIdQuery: func(c echo.Context, q *gorm.DB, entity T) error {
			tx := q.Delete(&entity)

			if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
				return ErrorResourceNotFound
			}

			if tx.Error != nil {
				return tx.Error
			}

			return nil
		},
	}
}

func (q *Queries[T]) ListByIdQuery(override func(c echo.Context, q *gorm.DB, id uint) (*T, error)) *Queries[T] {
	q.listByIdQuery = override

	return q
}

func (q *Queries[T]) ListAllQuery(override func(c echo.Context, q *gorm.DB) ([]T, error)) *Queries[T] {
	q.listAllQuery = override

	return q
}

func (q *Queries[T]) WriteByIdQuery(override func(c echo.Context, q *gorm.DB, entity *T, new any) error) *Queries[T] {
	q.writeByIdQuery = override

	return q
}

func (q *Queries[T]) DeleteByIdQuery(override func(c echo.Context, q *gorm.DB, entity T) error) *Queries[T] {
	q.deleteByIdQuery = override

	return q
}
