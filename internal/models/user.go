package models

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
	CompanyAddress    *string `json:"companyAddress"`
	PositionInCompany *string `json:"positionInCompany"`
	CompanyOkpo       *string `json:"companyOkpo"`
}
