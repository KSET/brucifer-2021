package repo

import (
	"github.com/google/uuid"

	"brucosijada.kset.org/app/models"
	"brucosijada.kset.org/app/providers/database"
)

type pageRepo struct {
}

type PageCreateOptions struct {
	Name     string
	Markdown string
}

type PageUpdateOptions struct {
	Model    *models.Page
	Name     string
	Markdown string
}

func Page() pageRepo {
	return pageRepo{}
}

func (r pageRepo) Create(data PageCreateOptions) (model *models.Page, err error) {
	db := database.DatabaseProvider().Client()

	model = &models.Page{
		Name:   data.Name,
		Markup: data.Markdown,
	}

	err = db.Create(model).Error

	return
}

func (r pageRepo) Update(data PageUpdateOptions) (err error) {
	db := database.DatabaseProvider().Client()

	data.Model.Name = data.Name
	data.Model.Markup = data.Markdown

	return db.Save(data.Model).Error
}

func (r pageRepo) List() (pages *[]models.Page, err error) {
	db := database.DatabaseProvider().Client()

	pages = &[]models.Page{}
	err = db.Model(&models.Page{}).Find(pages).Error

	return
}

func (r pageRepo) GetByUuid(id uuid.UUID) (model *models.Page, err error) {
	db := database.DatabaseProvider().Client()

	model = &models.Page{}
	err = db.Model(model).Where("uuid = ?", id.String()).First(model).Error

	return
}
