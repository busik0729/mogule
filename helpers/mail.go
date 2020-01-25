package helpers

import (
	"../config"

	"crypto/tls"

	"gopkg.in/gomail.v2"
)

type Mail struct {
	To      []string
	Message string
}

const FROM_EMAIL = "info@bstore.online"

func SendMail(mail Mail) {

	for _, emailAddress := range mail.To {
		m := gomail.NewMessage()
		m.SetHeader("From", FROM_EMAIL)
		m.SetHeader("To", emailAddress)
		// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
		m.SetHeader("Subject", "Info of BSTORE")
		m.SetBody("text/html", mail.Message)
		// m.Attach("/home/Alex/lolcat.jpg")

		d := initService()

		// Send the email to Bob, Cora and Dan.
		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}
	}

}

func initService() *gomail.Dialer {
	conf := config.GetSMTPConfig()

	d := gomail.NewDialer(conf.Host, conf.Port, conf.User, conf.Pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d
}
