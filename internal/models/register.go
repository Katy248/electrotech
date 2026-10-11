package models

type RegisterRequest struct {
	Email       string `binding:"required,email" json:"email"`
	Password    string `binding:"required,min=8" json:"password"`
	FirstName   string `binding:"required"       json:"first_name"`
	Surname     string `json:"surname"`
	LastName    string `json:"last_name"`
	PhoneNumber string `binding:"required"       json:"phone_number"`
}
