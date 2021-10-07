package repo

import (
	"brucosijada.kset.org/app/models"
	"brucosijada.kset.org/app/providers/database"
)

type userRepo struct {
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

func (u userRepo) GetById(id uint) models.User {
	db := database.DatabaseProvider().Client()

	var user models.User
	db.Where("id = ?", id).First(&user)

	return user
}
