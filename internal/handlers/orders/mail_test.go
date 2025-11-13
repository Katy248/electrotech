package orders

import (
	"database/sql"
	"electrotech/internal/repository/users"
	"os"
	"testing"
	"time"
)

func TestBuildMail(t *testing.T) {
	order := Order{
		ID:        21347328573298,
		UserID:    23578904367968,
		CreatedAt: time.Now(),
		Products: []Product{
			{
				ID:       "1",
				Name:     "Дилдо 20см",
				Price:    1000,
				Quantity: 120,
			},
			{
				ID:       "2",
				Name:     "Анальная пробка",
				Price:    2000,
				Quantity: 100,
			},
		},
	}
	user := users.User{
		ID:                124225,
		FirstName:         "Катерина",
		LastName:          "Владимировна",
		Email:             "katya@gmail.com",
		Surname:           "Васильева",
		PhoneNumber:       "+7436574398767695483",
		CompanyName:       sql.NullString{String: "ООО Наебалово и партнёры", Valid: true},
		CompanyAddress:    sql.NullString{String: "ул. Пушкина дом Колотушкина", Valid: true},
		CompanyInn:        sql.NullString{String: "1234567890I", Valid: true},
		CompanyOkpo:       sql.NullString{String: "1234567890O", Valid: true},
		PositionInCompany: sql.NullString{String: "Младший менеджер", Valid: true},
	}

	file, _ := os.Create("test.html")
	mail, err := buildMail(order, user)

	file.Write(mail)
	if err != nil {
		t.Errorf("Failed build mail: %s", err)
	}
}
