package users

import (
	"electrotech/internal/models"
	"electrotech/storage"
)

func ByEmail(email string) (*models.User, error) {
	var user models.User
	err := storage.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
func ByID(id int64) (*models.User, error) {
	var user models.User
	err := storage.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func InsertNew(u models.User) error {
	err := storage.DB.Create(&u).Error
	return err
}

func Update(u *models.User) error {
	err := storage.DB.Save(u).Error
	return err
}
