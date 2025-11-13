package orders

import (
	"bytes"
	"context"
	"electrotech/internal/email"
	"electrotech/internal/repository/users"
	"fmt"

	tmpl "html/template"

	"github.com/charmbracelet/log"
)

func sendEmail(order Order, userRepos *users.Queries) {

	if !email.IsEnabled() {
		return
	}

	user, err := userRepos.GetById(context.TODO(), order.UserID)
	if err != nil {
		log.Error("Failed getting user", "error", err, "userID", order.UserID)
		log.Error("Failed send email")
		return
	}

	err = email.SendInfo(buildMail(order, user))
	if err != nil {
		log.Error("Failed send email", "error", err)
	}

}

const mailTemplate = `
Subject: New order №{{ .ID }}

<h1> Заказ № {{ .ID }} </h1>

<p>Сумма заказа: {{ printf "%.2f" .Sum }} руб.</p>

Дата заказа: {{ .CreatedAt.Format "2006-01-02 15:04:05" }}

<h2> Товары </h2>

<ul>
{{ range .Products }}
	<li>
		<em>{{ .Name }}</em>
		({{ printf "%.2f" .Price }} руб.) x {{ .Quantity }} = {{ printf "%.2f" .Sum }} руб.</li>
{{ end }}
</ul>

## Заказчик

<p>{{ .User.Surname }} {{ .User.FirstName }} {{ .User.LastName }}</p>

<p>Почта: {{ .User.Email }}</p>

<p>Телефон: {{ .User.PhoneNumber }}</p>
`

func buildMail(order Order, user users.User) []byte {
	template, err := tmpl.New("new-order-mail").Parse(mailTemplate)
	if err != nil {
		log.Error("Failed parsing mail template", "error", err)
		panic(fmt.Sprintf("failed parse mail template: %s", err))
	}
	buff := bytes.Buffer{}
	template.Execute(&buff, order)
	return buff.Bytes()
}
