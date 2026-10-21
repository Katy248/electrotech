package order

import (
	"bytes"
	"electrotech/internal/models"
	_ "embed"
	"errors"
	"fmt"

	tmpl "html/template"

	"github.com/aymerick/douceur/inliner"
	"github.com/dustin/go-humanize"
)

func (s *Service) sendEmail(order models.Order) {
	buff, err := s.buildMail(order)
	if err != nil {
		s.logger.Error("Failed building mail", "error", err, "buffer", string(buff))

		return
	}

	err = s.emailService.SendInfo(buff, fmt.Sprintf("New Order #%d", order.ID))
	if err != nil {
		s.logger.Error("Failed send email", "error", err)

		return
	}

	s.logger.Info("Email sended")
}

var ErrOrderUserIsNil = errors.New("order user is nil")

//go:embed email.html
var EmailTemplate string

const FormatDecimals = 2

func (s *Service) buildMail(order models.Order) ([]byte, error) {
	if order.User == nil {
		return nil, ErrOrderUserIsNil
	}

	template, err := tmpl.New("new-order-mail").
		Funcs(tmpl.FuncMap{
			"humanizeComma": func(value float64) string {
				return humanize.CommafWithDigits(value, FormatDecimals)
			},
		}).
		Parse(EmailTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed parse template: %w", err)
	}

	buff := &bytes.Buffer{}

	err = template.Execute(buff, order)
	if err != nil {
		return buff.Bytes(), fmt.Errorf("failed execute template: %w", err)
	}

	inlined, err := inliner.Inline(buff.String())
	if err != nil {
		return nil, fmt.Errorf("failed inline styles: %w", err)
	}

	return []byte(inlined), nil
}
