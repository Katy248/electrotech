package orders

import (
	"bytes"
	"electrotech/internal/models"
	_ "embed"
	"errors"
	"fmt"

	tmpl "html/template"

	"charm.land/log/v2"
	"github.com/aymerick/douceur/inliner"
)

func (h *Handler) sendEmail(order models.Order) {
	buff, err := h.buildMail(order)
	if err != nil {
		h.logger.Error("Failed building mail", "error", err, "buffer", string(buff))

		return
	}

	err = h.EmailService.SendInfo(buff, fmt.Sprintf("New Order #%d", order.ID))
	if err != nil {
		h.logger.Error("Failed send email", "error", err)
	}
}

var ErrOrderUserIsNil = errors.New("order user is nil")

//go:embed email.html
var EmailTemplate string

func (h *Handler) buildMail(order models.Order) ([]byte, error) {
	if order.User == nil {
		return nil, ErrOrderUserIsNil
	}

	template, err := tmpl.New("new-order-mail").Parse(EmailTemplate)
	if err != nil {
		h.logger.Error("Failed parsing mail template", "error", err)

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
