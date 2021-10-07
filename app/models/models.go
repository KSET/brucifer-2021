package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uint      `json:"-" gorm:"primaryKey"`
	UUID      uuid.UUID `json:"_id" gorm:"index;autoIncrement:false"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// BeforeCreate will set Base struct before every insert
func (base *Base) BeforeCreate(tx *gorm.DB) error {
	// uuid.New() creates a new random UUID or panics.
	base.UUID = uuid.New()

	return nil
}

func (base Base) Exists() bool {
	return base.ID != 0
}
