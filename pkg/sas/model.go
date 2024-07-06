package sas

import (
	"gorm.io/gorm"
	"time"
)

// DefaultModel is a normal db model, with a uint primary key.
type DefaultModel struct {
	ID      uint `gorm:"primaryKey"`
	Created time.Time
	Updated time.Time
}

func (m *DefaultModel) BeforeCreate(tx *gorm.DB) error {
	utcTime := time.Now().UTC()
	m.Created = utcTime
	m.Updated = utcTime

	return nil
}

func (m *DefaultModel) BeforeUpdate(tx *gorm.DB) error {
	utcTime := time.Now().UTC()
	m.Updated = utcTime

	return nil
}
