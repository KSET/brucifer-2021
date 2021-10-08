package repo

import (
	"gorm.io/gorm"

	"brucosijada.kset.org/app/models"
	"brucosijada.kset.org/app/providers/database"
	"brucosijada.kset.org/app/providers/hash"
)

type userRepo struct {
}

type UserCreateModel struct {
	Email      string
	Identity   string
	Password   string
	Invitation *models.UserInvitation
}

func User() userRepo {
	return userRepo{}
}

func (u userRepo) GetByUsername(username string) models.User {
	db := database.DatabaseProvider().Client()

	var user models.User
	db.Where("username = ?", username).First(&user)

	return user
}

func (u userRepo) GetByIdentity(identity string) models.User {
	db := database.DatabaseProvider().Client()

	var user models.User
	db.Where("username = ? or email = ?", identity, identity).First(&user)

	return user
}

func (u userRepo) GetById(id uint) models.User {
	db := database.DatabaseProvider().Client()

	var user models.User
	db.Where("id = ?", id).First(&user)

	return user
}

func (u userRepo) List() (users []models.User) {
	db := database.DatabaseProvider().Client()

	db.Find(&users)

	return users
}

func (u userRepo) Create(model UserCreateModel) (
	user models.User,
	err error,
) {
	err = database.DatabaseProvider().Client().Transaction(
		func(tx *gorm.DB) (err error) {
			user = models.User{
				Email:    model.Email,
				Username: model.Identity,
				Password: hash.HashProvider().HashPassword(model.Password),
			}

			if err = tx.Save(&user).Error; err != nil {
				return
			}

			if model.Invitation != nil && model.Invitation.Exists() {
				model.Invitation.UsedByID = user.ID

				if err = tx.Save(model.Invitation).Error; err != nil {
					return
				}
			}

			return nil
		},
	)

	return
}
