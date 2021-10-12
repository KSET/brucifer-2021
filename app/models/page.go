package models

import (
	"gorm.io/gorm"

	"brucosijada.kset.org/app/providers/markdown"
)

type Page struct {
	Base
	Name     string `json:"name" gorm:"unique;index;size:127"`
	Markup   string `json:"contents"`
	Rendered string `json:"rendered"`
}

func (p *Page) BeforeSave(tx *gorm.DB) (err error) {
	p.Rendered = markdown.MarkdownProvider().Render(p.Markup)

	return
}
