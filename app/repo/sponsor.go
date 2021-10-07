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

type SponsorItemLogo struct {
	Width uint   `json:"width"`
	Url   string `json:"url"`
}

type SponsorItem struct {
	Id    uuid.UUID         `json:"id"`
	Name  string            `json:"name"`
	Logo  []SponsorItemLogo `json:"logo"`
}

func Sponsor() sponsorRepo {
	return sponsorRepo{}
}

func (r sponsorRepo) Create(name string, logo *multipart.FileHeader, uploaderId uint) (
	sponsor *models.Sponsor,
	err error,
) {
	image, err := Image().Create(logo, uploaderId)

	if err != nil {
		return nil, err
	}

	db := database.DatabaseProvider().Client()

	err = db.Transaction(
		func(tx *gorm.DB) error {
			var orderSponsor models.Sponsor
			tx.Order("\"order\" desc").First(&orderSponsor)

			sponsor = &models.Sponsor{
				Name:  name,
				Image: *image,
				Order: orderSponsor.Order + 1,
			}

			return tx.Save(sponsor).Error
		},
	)

	return sponsor, err
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
			Id:    sponsor.UUID,
			Name:  sponsor.Name,
			Logo:  logos,
		}
	}

	return &items, nil
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
