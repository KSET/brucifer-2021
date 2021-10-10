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

type artistRepo struct {
}

type ArtistCreateOptions struct {
	Name     string
	Logo     *multipart.FileHeader
	Uploader models.User
}

type ArtistUpdateOptions struct {
	Model    *models.Artist
	Name     string
	Logo     *multipart.FileHeader
	Uploader models.User
}

type ArtistItemLogo struct {
	Width uint   `json:"width"`
	Url   string `json:"url"`
}

type ArtistItem struct {
	Id   uuid.UUID        `json:"id"`
	Name string           `json:"name"`
	Logo []ArtistItemLogo `json:"logo"`
}

func Artist() artistRepo {
	return artistRepo{}
}

func (r artistRepo) Create(data ArtistCreateOptions) (
	artist *models.Artist,
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
			var orderArtist models.Artist
			tx.Order("\"order\" desc").First(&orderArtist)

			artist = &models.Artist{
				Name:  data.Name,
				Image: *image,
				Order: orderArtist.Order + 1,
			}

			return tx.Save(artist).Error
		},
	)

	return artist, err
}

func (r artistRepo) Update(data ArtistUpdateOptions) error {
	db := database.DatabaseProvider().Client()

	return db.Transaction(
		func(tx *gorm.DB) (err error) {
			if data.Logo != nil {
				if err = tx.Delete(&data.Logo).Error; err != nil {
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

			return tx.Save(data.Model).Error
		},
	)
}

func (r artistRepo) ListSimple() (*[]ArtistItem, error) {
	var artists []*models.Artist
	err := database.DatabaseProvider().Client().Order("\"order\" ASC").Preload("Image.Variations").Find(&artists).Error
	if err != nil {
		return nil, err
	}

	items := make([]ArtistItem, len(artists))
	for i, artist := range artists {
		items[i] = *r.ToArtistItem(artist)
	}

	return &items, nil
}

func (r artistRepo) GetByUuidSimple(id uuid.UUID) (artistItem *ArtistItem, err error) {
	var artist models.Artist
	q := database.DatabaseProvider().Client().Where(
		"uuid = ?",
		id.String(),
	).Order(
		"\"order\" ASC",
	).Preload(
		"Image.Variations",
	).Find(&artist)

	err = q.Error
	if err != nil {
		return
	}

	if q.RowsAffected == 0 {
		err = errors.New("Item not found")
		return
	}

	artistItem = r.ToArtistItem(&artist)
	return
}

func (r artistRepo) Swap(swap, with uuid.UUID) error {
	db := database.DatabaseProvider().Client()

	err := db.Transaction(
		func(tx *gorm.DB) error {
			var artists []models.Artist
			err := tx.Where(
				"uuid in (?)",
				[]uuid.UUID{
					swap,
					with,
				},
			).Find(&artists).Error

			if err != nil {
				return err
			}

			if len(artists) != 2 {
				return errors.New("Artist does not exist")
			}

			s, w := artists[0], artists[1]

			if err := db.Model(&models.Artist{}).Where("uuid = ?", s.UUID).Update("order", w.Order).Error; err != nil {
				return err
			}

			if err := db.Model(&models.Artist{}).Where("uuid = ?", w.UUID).Update("order", s.Order).Error; err != nil {
				return err
			}

			return nil
		},
	)

	return err
}

func (r artistRepo) ToArtistItem(model *models.Artist) (item *ArtistItem) {
	logos := make([]ArtistItemLogo, len(model.Image.Variations))
	for i, variation := range model.Image.Variations {
		logos[i] = ArtistItemLogo{
			Width: variation.Width,
			Url:   fmt.Sprintf("/i/%s", variation.UUID.String()),
		}
	}

	return &ArtistItem{
		Id:   model.UUID,
		Name: model.Name,
		Logo: logos,
	}
}
