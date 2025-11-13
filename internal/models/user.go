package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                int64   `json:"id"`
	FirstName         string  `json:"firstName"`
	Surname           string  `json:"surname"`
	LastName          string  `json:"lastName"`
	Email             string  `json:"email"`
	PhoneNumber       string  `json:"phoneNumber"`
	PasswordHash      string  `json:"-"`
	CompanyName       *string `json:"companyName"`
	CompanyInn        *string `json:"companyInn"`
	CompanyOkpo       *string `json:"companyOkpo"`
	CompanyAddress    *string `json:"companyAddress"`
	PositionInCompany *string `json:"positionInCompany"`
}

func (u *User) UpdatePassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed hash password: %s", err)
	}
	u.PasswordHash = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) CompanyData() *CompanyData {
	return &CompanyData{
		Name:     *u.CompanyName,
		INN:      *u.CompanyInn,
		OKPO:     *u.CompanyOkpo,
		Address:  *u.CompanyAddress,
		Position: *u.PositionInCompany,
	}
}

type CompanyData struct {
	Name     string
	INN      string
	OKPO     string
	Address  string
	Position string
}

func (c *CompanyData) DataFilled() bool {
	return c.Name != "" && c.INN != "" && c.OKPO != "" && c.Address != "" && c.Position != ""
}
