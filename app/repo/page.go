package repo

import (
	"mime/multipart"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"brucosijada.kset.org/app/models"
	"brucosijada.kset.org/app/providers/database"
)

type pageRepo struct {
}

type PageCreateOptions struct {
	Name             string
	Markdown         string
	Uploader         models.User
	Background       *multipart.FileHeader
	BackgroundMobile *multipart.FileHeader
}

type PageUpdateOptions struct {
	Model             *models.Page
	Name              string
	Markdown          string
	Uploader          models.User
	Background        *multipart.FileHeader
	BackgroundM       *models.Image
	BackgroundMobile  *multipart.FileHeader
	BackgroundMobileM *models.Image
}

type PageItem struct {
	ID               uuid.UUID  `json:"id"`
	Name             string     `json:"name"`
	Markup           string     `json:"contents"`
	Rendered         string     `json:"rendered"`
	Background       *ImageItem `json:"background"`
	BackgroundMobile *ImageItem `json:"backgroundMobile"`
}

func Page() pageRepo {
	return pageRepo{}
}

func (r pageRepo) Create(data PageCreateOptions) (model *models.Page, err error) {
	db := database.DatabaseProvider().Client()

	var bg, bgMobile *models.Image

	if data.Background != nil {
		bg, err = Image().Create(
			ImageCreateOptions{
				File:     data.Background,
				Uploader: data.Uploader,
			},
		)
		if err != nil {
			return
		}
	}

	if data.BackgroundMobile != nil {
		bgMobile, err = Image().Create(
			ImageCreateOptions{
				File:     data.BackgroundMobile,
				Uploader: data.Uploader,
			},
		)
		if err != nil {
			return
		}
	}

	model = &models.Page{
		Name:             data.Name,
		Markup:           data.Markdown,
		Background:       bg,
		BackgroundMobile: bgMobile,
	}

	err = db.Create(model).Error

	return
}

func (r pageRepo) Update(data PageUpdateOptions) (err error) {
	var bg, bgMobile *models.Image

	if data.Background != nil {
		bg, err = Image().Create(
			ImageCreateOptions{
				File:     data.Background,
				Uploader: data.Uploader,
			},
		)
		if err != nil {
			return
		}
		data.BackgroundM = bg
	}

	if data.BackgroundMobile != nil {
		bgMobile, err = Image().Create(
			ImageCreateOptions{
				File:     data.BackgroundMobile,
				Uploader: data.Uploader,
			},
		)
		if err != nil {
			return
		}
		data.BackgroundMobileM = bgMobile
	}

	data.Model.Name = data.Name
	data.Model.Markup = data.Markdown

	bgOld := data.Model.Background
	data.Model.Background = data.BackgroundM

	bgMobileOld := data.Model.BackgroundMobile
	data.Model.BackgroundMobile = data.BackgroundMobileM

	return database.DatabaseProvider().Client().Transaction(
		func(tx *gorm.DB) (err error) {
			if bgOld != nil && bgOld != data.Model.Background {
				if err = tx.Delete(bgOld).Error; err != nil {
					return
				}
			}

			if bgMobileOld != nil && bgMobileOld != data.Model.BackgroundMobile {
				if err = tx.Delete(bgMobileOld).Error; err != nil {
					return
				}
			}

			return tx.Save(data.Model).Error
		},
	)
}

func (r pageRepo) ListSimple() (pageItems *[]PageItem, err error) {
	db := database.DatabaseProvider().Client()

	var pages []models.Page
	err = db.Model(&models.Page{}).Preload("Background.Variations").Preload("BackgroundMobile.Variations").Find(&pages).Error
	if err != nil {
		return
	}

	items := make([]PageItem, 0)
	for _, page := range pages {
		items = append(items, r.ToPageItem(&page))
	}
	pageItems = &items

	return
}

func (r pageRepo) GetByUuid(id uuid.UUID) (model *models.Page, err error) {
	db := database.DatabaseProvider().Client()

	model = &models.Page{}
	err = db.Model(model).Preload("Background.Variations").Preload("BackgroundMobile.Variations").Where(
		"uuid = ?",
		id.String(),
	).First(model).Error

	return
}

func (r pageRepo) ToPageItem(model *models.Page) (item PageItem) {
	bg := Image().ToImageItem(model.Background)
	bgMobile := Image().ToImageItem(model.BackgroundMobile)

	data := PageItem{
		ID:               model.UUID,
		Name:             model.Name,
		Markup:           model.Markup,
		Rendered:         model.Rendered,
		Background:       bg,
		BackgroundMobile: bgMobile,
	}

	return data
}
