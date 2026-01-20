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

// Receives password string and sets PasswordHash to hash of input
func (u *User) SetPassword(password string) error {
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
	data := &CompanyData{
		Name:     strValueOrEmpty(u.CompanyName),
		INN:      strValueOrEmpty(u.CompanyInn),
		Address:  strValueOrEmpty(u.CompanyAddress),
		OKPO:     strValueOrEmpty(u.CompanyOkpo),
		Position: strValueOrEmpty(u.PositionInCompany),
	}
	data.AllRequiredFields = data.DataFilled()
	return data
}

// Utility func to work with nullable representation of sql strings
func strValueOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

type CompanyData struct {
	Name     string `json:"companyName"`
	INN      string `json:"companyInn"`
	OKPO     string `json:"companyOkpo"`
	Address  string `json:"companyAddress"`
	Position string `json:"positionInCompany"`

	AllRequiredFields bool `json:"allRequiredFields"`
}

func (c *CompanyData) DataFilled() bool {
	return c.Name != "" && c.INN != "" && c.OKPO != "" && c.Address != "" && c.Position != ""
}
