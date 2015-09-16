package mail_test

import (
	"./" // mail package
	"bytes"
	"fmt"
	"html/template"
	"log"
	"testing"

	. "github.com/etsy/mixer/config"
)

func TestMail(*testing.T) {
	err := Config.Load()
	if err != nil {
		log.Fatal("error reading or parsing config:", err)
	}

	m, _ := mail.NewMail()

	mdir := Config.GetRootDir()
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/mail/match_email.tmpl", mdir))
	if err != nil {
		log.Fatal(err)
	}

	emailData := struct {
		DirectoryURL string
		Person1LDAP  string
		Person1Name  string
		Person2LDAP  string
		Person2Name  string
		MixerName    string
		ServerUrl    string
	}{
		Config.Staff.DirectoryUrl,
		"thing1ldap",
		"thing1",
		"thing2ldap",
		"thing2",
		"Managers",
		Config.Server.Url,
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, emailData)
	if err != nil {
		log.Fatal(err)
	}

	subject := "Managers Mixer"

	var to = []string{fmt.Sprintf("%s@%s", Config.Mail.AdminUsername, Config.Mail.Domain)}
	to = append(to, fmt.Sprintf("%s+secondto@%s", Config.Mail.AdminUsername, Config.Mail.Domain))

	m.Mail(buffer.Bytes(), subject, to)
	fmt.Println("test")
}
