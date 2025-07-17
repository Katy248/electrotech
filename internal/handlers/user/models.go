package user

type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	FirstName   string `json:"first_name" binding:"required"`
	Surname     string `json:"surname" binding:"required"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token     string `json:"token"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	Surname   string `json:"surname"`
}
