package repo

import (
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"brucosijada.kset.org/app/models"
	"brucosijada.kset.org/app/providers/database"
)

type sponsorRepo struct {
}

type SponsorCreateOptions struct {
	Name     string
	Link     string
	Logo     *multipart.FileHeader
	Uploader models.User
}

type SponsorUpdateOptions struct {
	Model    *models.Sponsor
	Name     string
	Link     string
	Logo     *multipart.FileHeader
	Uploader models.User
}

type SponsorItemLogo struct {
	Width uint   `json:"width"`
	Url   string `json:"url"`
}

type SponsorItem struct {
	Id   uuid.UUID         `json:"id"`
	Name string            `json:"name"`
	Url  string            `json:"url"`
	Logo []SponsorItemLogo `json:"logo"`
}

func Sponsor() sponsorRepo {
	return sponsorRepo{}
}

func (r sponsorRepo) Create(data SponsorCreateOptions) (
	sponsor *models.Sponsor,
	err error,
) {
	image, err := Image().Create(
		ImageCreateOptions{
			File:     data.Logo,
			Uploader: data.Uploader,
		},
	)

	if err != nil {
		return nil, err
	}

	db := database.DatabaseProvider().Client()

	err = db.Transaction(
		func(tx *gorm.DB) error {
			var orderSponsor models.Sponsor
			tx.Order("\"order\" desc").First(&orderSponsor)

			sponsor = &models.Sponsor{
				Name:  data.Name,
				Link:  data.Link,
				Image: *image,
				Order: orderSponsor.Order + 1,
			}

			return tx.Save(sponsor).Error
		},
	)

	return sponsor, err
}

func (r sponsorRepo) Update(data SponsorUpdateOptions) error {
	db := database.DatabaseProvider().Client()

	return db.Transaction(
		func(tx *gorm.DB) (err error) {
			if data.Logo != nil {
				if err = tx.Delete(&data.Model.Image).Error; err != nil {
					return
				}

				var image *models.Image
				if image, err = Image().Create(
					ImageCreateOptions{
						File:     data.Logo,
						Uploader: data.Uploader,
					},
				); err != nil {
					return
				}

				data.Model.Image = *image
			}

			if err != nil {
				return
			}

			data.Model.Name = data.Name
			data.Model.Link = data.Link

			return tx.Save(data.Model).Error
		},
	)
}

func (r sponsorRepo) ListSimple() (*[]SponsorItem, error) {
	var sponsors []*models.Sponsor
	err := database.DatabaseProvider().Client().Order("\"order\" ASC").Preload("Image.Variations").Find(&sponsors).Error
	if err != nil {
		return nil, err
	}

	items := make([]SponsorItem, len(sponsors))
	for i, sponsor := range sponsors {
		logos := make([]SponsorItemLogo, len(sponsor.Image.Variations))
		for i, variation := range sponsor.Image.Variations {
			logos[i] = SponsorItemLogo{
				Width: variation.Width,
				Url:   fmt.Sprintf("/i/%s", variation.UUID.String()),
			}
		}

		items[i] = SponsorItem{
			Id:   sponsor.UUID,
			Name: sponsor.Name,
			Url:  sponsor.Link,
			Logo: logos,
		}
	}

	return &items, nil
}

func (r sponsorRepo) GetByUuidSimple(id uuid.UUID) (sponsorItem *SponsorItem, err error) {
	var sponsor models.Sponsor
	q := database.DatabaseProvider().Client().Where(
		"uuid = ?",
		id.String(),
	).Order(
		"\"order\" ASC",
	).Preload(
		"Image.Variations",
	).Find(&sponsor)

	err = q.Error
	if err != nil {
		return
	}

	if q.RowsAffected == 0 {
		err = errors.New("Item not found")
		return
	}

	sponsorItem = r.ToSponsorItem(&sponsor)
	return
}

func (r sponsorRepo) Swap(swap, with uuid.UUID) error {
	db := database.DatabaseProvider().Client()

	err := db.Transaction(
		func(tx *gorm.DB) error {
			var sponsors []models.Sponsor
			err := tx.Where(
				"uuid in (?)",
				[]uuid.UUID{
					swap,
					with,
				},
			).Find(&sponsors).Error

			if err != nil {
				return err
			}

			if len(sponsors) != 2 {
				return errors.New("Sponsor does not exist")
			}

			s, w := sponsors[0], sponsors[1]

			if err := db.Model(&models.Sponsor{}).Where("uuid = ?", s.UUID).Update("order", w.Order).Error; err != nil {
				return err
			}

			if err := db.Model(&models.Sponsor{}).Where("uuid = ?", w.UUID).Update("order", s.Order).Error; err != nil {
				return err
			}

			return nil
		},
	)

	return err
}

func (r sponsorRepo) ToSponsorItem(model *models.Sponsor) (item *SponsorItem) {
	logos := make([]SponsorItemLogo, len(model.Image.Variations))
	for i, variation := range model.Image.Variations {
		logos[i] = SponsorItemLogo{
			Width: variation.Width,
			Url:   fmt.Sprintf("/i/%s", variation.UUID.String()),
		}
	}

	return &SponsorItem{
		Id:   model.UUID,
		Name: model.Name,
		Url:  model.Link,
		Logo: logos,
	}
}
