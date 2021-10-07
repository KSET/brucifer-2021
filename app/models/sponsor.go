package models

import (
	"gorm.io/gorm"
)

type Sponsor struct {
	Base
	Name    string `json:"name"`
	Image   Image  `json:"image" gorm:"polymorphic:Owner"`
	Order   int    `json:"order" gorm:"index"`
}

// AfterDelete hook defined for cascade delete
func (r *Sponsor) AfterDelete(tx *gorm.DB) error {
	err := tx.Delete(&r.Image).Error

	return err
}