package mail

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/etsy/mixer/config"
)

type EmailUser struct {
	Username    string
	Password    string
	EmailServer string
	Port        int
}

func NewMail() (*EmailUser, error) {

	err := config.Config.Load()
	if err != nil {
		fmt.Println("error reading or parsing config config:", err)
	}

	emailUser := &EmailUser{config.Config.Mail.EmailAddress, config.Config.Mail.Key, config.Config.Mail.SMTP_Host, config.Config.Mail.Port}
	return emailUser, nil
}

func (eu *EmailUser) Mail(msg []byte, subject string, to []string) {

	auth := smtp.PlainAuth(
		"",
		eu.Username,
		eu.Password,
		eu.EmailServer,
	)

	address := fmt.Sprintf("%v:%v", eu.EmailServer, eu.Port)

	header := make(map[string]string)
	header["From"] = eu.Username
	header["To"] = strings.Join(to, ",")
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString(msg)

	err := smtp.SendMail(
		address,
		auth,
		eu.Username,
		to,
		[]byte(message),
	)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
