package models

type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	Surname      string `json:"surname"`
	LastName     string `json:"last_name"`
	PhoneNumber  string `json:"phone_number"`
}
