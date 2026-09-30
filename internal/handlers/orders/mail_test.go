package orders

import (
	"electrotech/internal/models"
	"os"
	"testing"
	"time"

	"charm.land/log/v2"
)

func TestBuildMail(t *testing.T) {
	t.Parallel()

	order := models.Order{ //nolint:exhaustruct_v5
		ID:           21347328573298,
		UserID:       23578904367968,
		CreationDate: time.Now(),
		OrderProducts: []*models.OrderProduct{
			{ //nolint:exhaustruct_v5
				ID:           1,
				ProductID:    "1",
				ProductName:  "Дилдо 20см",
				ProductPrice: 1000,
				Quantity:     120,
			},
			{ //nolint:exhaustruct_v5
				ProductID:    "2",
				ProductName:  "Анальная пробка",
				ProductPrice: 2000,
				Quantity:     100,
			},
		},
	}
	user := models.User{ //nolint:exhaustruct_v5
		ID:          124225,
		FirstName:   "Катерина",
		LastName:    "Владимировна",
		Email:       "katya@gmail.com",
		Surname:     "Васильева",
		PhoneNumber: "+7436574398767695483",
	}
	{
		name := "ООО Наебалово и партнёры"
		user.CompanyName = &name
		addr := "ул. Пушкина дом Колотушкина"
		user.CompanyAddress = &addr
		inn := "1234567890I"
		user.CompanyInn = &inn
		okpo := "1234567890O"
		user.CompanyOkpo = &okpo
		pos := "Младший менеджер"
		user.PositionInCompany = &pos
	}

	order.User = &user

	file, _ := os.Create("test.html")

	handler := &Handler{ //nolint:exhaustruct_v5
		logger: log.New(os.Stdout),
	}

	mail, err := handler.buildMail(order)
	if err != nil {
		t.Errorf("Failed build mail: %s", err)
	}

	_, err = file.Write(mail)
	if err != nil {
		t.Errorf("Failed write mail: %s", err)
	}
}
