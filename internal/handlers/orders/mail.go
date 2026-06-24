package orders

import (
	"bytes"
	"electrotech/internal/email"
	"electrotech/internal/models"
	_ "embed"
	"errors"
	"fmt"

	tmpl "html/template"

	"github.com/aymerick/douceur/inliner"
	"github.com/charmbracelet/log"
)

func sendEmail(order models.Order) {
	if !email.IsEnabled() {
		return
	}

	buff, err := buildMail(order)
	if err != nil {
		log.Error("Failed building mail", "error", err, "buffer", string(buff))

		return
	}

	err = email.SendInfo(buff, fmt.Sprintf("New Order #%d", order.ID))
	if err != nil {
		log.Error("Failed send email", "error", err)
	}
}

var ErrOrderUserIsNil = errors.New("order user is nil")

//go:embed email.html
var EmailTemplate string

func buildMail(order models.Order) ([]byte, error) {
	if order.User == nil {
		return nil, ErrOrderUserIsNil
	}

	template, err := tmpl.New("new-order-mail").Parse(EmailTemplate)
	if err != nil {
		log.Error("Failed parsing mail template", "error", err)

		return nil, fmt.Errorf("failed parse template: %w", err)
	}

	buff := &bytes.Buffer{}

	err = template.Execute(buff, order)
	if err != nil {
		log.Error("Failed executing mail template", "error", err, "order", order)

		return buff.Bytes(), fmt.Errorf("failed execute template: %w", err)
	}

	inlined, err := inliner.Inline(buff.String())
	if err != nil {
		return nil, fmt.Errorf("failed inline styles: %w", err)
	}

	return []byte(inlined), nil
}
