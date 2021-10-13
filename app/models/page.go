package models

import (
	"gorm.io/gorm"

	"brucosijada.kset.org/app/providers/markdown"
)

type Page struct {
	Base
	Name             string `json:"name" gorm:"unique;index;size:127"`
	Markup           string `json:"contents"`
	Rendered         string `json:"rendered"`
	Background       *Image `json:"background" gorm:"polymorphic:Owner"`
	BackgroundMobile *Image `json:"backgroundMobile" gorm:"polymorphic:Owner;polymorphicValue:pages_mobile"`
}

func (p *Page) BeforeSave(tx *gorm.DB) (err error) {
	p.Rendered = markdown.MarkdownProvider().Render(p.Markup)

	return
}

// AfterDelete hook defined for cascade delete
func (p *Page) AfterDelete(tx *gorm.DB) error {
	if p.Background != nil {
		err := tx.Delete(p.Background).Error
		if err != nil {
			return err
		}
	}

	if p.BackgroundMobile != nil {
		err := tx.Delete(p.BackgroundMobile).Error
		if err != nil {
			return err
		}
	}

	return nil
}
