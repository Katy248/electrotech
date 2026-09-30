package users

import (
	"electrotech/internal/models"
	"strings"

	"charm.land/log/v2"
	"gorm.io/gorm"
)

type Repo struct {
	DB     *gorm.DB
	logger *log.Logger
}

func NewRepo(db *gorm.DB, logger *log.Logger) *Repo {
	return &Repo{DB: db, logger: logger}
}

func (r *Repo) ByEmail(email string) (*models.User, error) {
	email = strings.ToLower(email)

	var user models.User

	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repo) ByID(id int64) (*models.User, error) {
	var user models.User

	err := r.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repo) InsertNew(u *models.User) error {
	normalizeEmail(u)
	err := r.DB.Create(&u).Error

	return err
}

func (r *Repo) Update(u *models.User) error {
	normalizeEmail(u)
	err := r.DB.Save(u).Error

	return err
}

func normalizeEmail(u *models.User) {
	u.Email = strings.ToLower(u.Email)
}
