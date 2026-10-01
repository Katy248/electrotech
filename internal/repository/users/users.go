package users

import (
	"context"
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

func (r *Repo) ByEmail(ctx context.Context, email string) (*models.User, error) {
	email = strings.ToLower(email)

	var user models.User

	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repo) ByID(ctx context.Context, id int64) (*models.User, error) {
	var user models.User

	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repo) InsertNew(ctx context.Context, u *models.User) error {
	normalizeEmail(u)
	err := r.DB.WithContext(ctx).Create(&u).Error

	return err
}

func (r *Repo) Update(ctx context.Context, u *models.User) error {
	normalizeEmail(u)
	err := r.DB.Save(u).Error

	return err
}

func normalizeEmail(u *models.User) {
	u.Email = strings.ToLower(u.Email)
}
