package orders

import (
	"bytes"
	"electrotech/internal/email"
	"electrotech/internal/models"
	"electrotech/internal/users"
	_ "embed"
	"fmt"

	tmpl "html/template"

	"github.com/aymerick/douceur/inliner"
	"github.com/charmbracelet/log"
)

func sendEmail(order Order) {

	if !email.IsEnabled() {
		return
	}

	user, err := users.ByID(order.UserID)
	if err != nil {
		log.Error("Failed getting user", "error", err, "userID", order.UserID)
		log.Error("Failed send email")
		return
	}
	buff, err := buildMail(order, *user)
	if err != nil {
		log.Error("Failed building mail", "error", err, "buffer", buff)
		return
	}

	err = email.SendInfo(buff, fmt.Sprintf("New Order #%d", order.ID))
	if err != nil {
		log.Error("Failed send email", "error", err)
	}
}

//go:embed email.html
var EmailTemplate string

func buildMail(order Order, user models.User) ([]byte, error) {
	template, err := tmpl.New("new-order-mail").Parse(EmailTemplate)

	if err != nil {
		log.Error("Failed parsing mail template", "error", err)
		return nil, fmt.Errorf("failed parse template: %s", err)
	}
	data := struct {
		Order Order
		User  models.User
	}{
		Order: order,
		User:  user,
	}
	buff := &bytes.Buffer{}
	err = template.Execute(buff, data)
	if err != nil {
		log.Error("Failed executing mail template", "error", err)
		return buff.Bytes(), fmt.Errorf("failed execute template: %s", err)
	}

	inlined, err := inliner.Inline(buff.String())
	if err != nil {
		return nil, fmt.Errorf("failed inline styles: %s", err)
	}

	return []byte(inlined), nil
}
