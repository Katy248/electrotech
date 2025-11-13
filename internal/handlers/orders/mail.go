package orders

import (
	"bytes"
	"context"
	"electrotech/internal/email"
	"electrotech/internal/repository/users"
	"fmt"
	"io"

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
	buff := bytes.Buffer{}
	buildMail(order, user, &buff)

	err = email.SendInfo(buff.Bytes(), fmt.Sprintf("New Order #%d", order.ID))
	if err != nil {
		log.Error("Failed send email", "error", err)
	}
}

const mailTemplate = `
<!DOCTYPE html>
<html lang="ru">
<head>
	<meta charset="UTF-8">
</head>
<style>
	body {
		font-family: sans-serif;
		background-color: #f5f5f5;
		font-size: 16px;
	}
	main {
		margin: 1rem;
		padding: 2rem;
		background-color: #fff;
		border-radius: 10px;
		shadow: 0 0 10px rgba(0, 0, 0, 0.1);
	}

	h1, h2, h3 {
		color: #0081b6;
	}
	h1 {
		font-size: 24px;
	}
	h2{
		font-size: 18px;
	}
	h3{
		font-size: 16px;
	}
</style>
<body>

<main>
<h1> Новый заказ № {{ .Order.ID }} </h1>

<p>Сумма заказа: {{ .Order.FormatSum "руб." }} </p>

Дата заказа: {{ .Order.CreatedAt.Format "2006-01-02 15:04:05" }}

<h2> Товары </h2>

<ul>
{{ range .Order.Products }}
	<li>
		<em>{{ .Name }}</em>
		({{ printf "%.2f" .Price }} руб.) x {{ .Quantity }} = {{ .FormatSum "руб." }}</li>
{{ end }}
</ul>

<h2> Заказчик </h2>

<h3><q>{{ .User.CompanyName.Value }}</q></h3>


<p>
Адрес: {{ .User.CompanyAddress.Value }}
</p>
<p>
ИНН: {{ .User.CompanyInn.Value }}
</p>
<p>
ОКПО: {{ .User.CompanyOkpo.Value }}
</p>

<h3>Контакты</h3>

<p>{{ .User.Surname }} {{ .User.FirstName }} {{ .User.LastName }} - {{ .User.PositionInCompany.Value }}</p>


<p>Почта: <a href="mailto:{{ .User.Email }}">{{ .User.Email }}</a></p>

<p>Телефон: <a href="tel:{{ .User.PhoneNumber }}">{{ .User.PhoneNumber }}</a></p>
</main>
</body>
</html>
`

func buildMail(order Order, user users.User, w io.Writer) error {
	template, err := tmpl.New("new-order-mail").Parse(mailTemplate)
	if err != nil {
		log.Error("Failed parsing mail template", "error", err)
		return fmt.Errorf("failed parse template: %s", err)
	}
	data := struct {
		Order Order
		User  users.User
	}{
		Order: order,
		User:  user,
	}
	err = template.Execute(w, data)
	if err != nil {
		log.Error("Failed executing mail template", "error", err)
		return fmt.Errorf("failed execute template: %s", err)
	}
	return nil
}
