package users

import (
	"electrotech/internal/models"
	"electrotech/storage"
	"strings"
)

func ByEmail(email string) (*models.User, error) {
	email = strings.ToLower(email)
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

func normalizeEmail(u *models.User) {
	u.Email = strings.ToLower(u.Email)
}

func InsertNew(u *models.User) error {
	normalizeEmail(u)
	err := storage.DB.Create(&u).Error
	return err
}

func Update(u *models.User) error {
	normalizeEmail(u)
	err := storage.DB.Save(u).Error
	return err
}
