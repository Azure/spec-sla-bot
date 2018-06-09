package email

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"

	"github.com/Azure/spec-sla-bot/template"
	gomail "gopkg.in/gomail.v2"
)

//SendEmailToAssignee sends an email to a list of users
func SendEmailToAssignee() error {
	template.GenerateTemplate()
	b, err := ioutil.ReadFile("finalTemplate.html")
	if err != nil {
		fmt.Print(err)
		return err
	}
	str := string(b) // convert content to a 'string'
	//fmt.Print(str)
	m := gomail.NewMessage()
	//Get connection string from azure
	emailUrl := os.Getenv("CUSTOMCONNSTR_EMAIL_URL")

	//parse connection string url
	parsed, err := url.Parse(emailUrl)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	port, _ := strconv.Atoi(parsed.Port())
	user := parsed.User.Username()
	password, _ := parsed.User.Password()
	email := user + "@" + parsed.Hostname()
	m.SetHeader("From", email)
	m.SetHeader("To", "t-jaelli@microsoft.com")
	m.SetHeader("Subject", "TEST")
	m.SetBody("text/html", str)

	//Send the email
	d := gomail.NewDialer("smtp.office365.com", port, email, password)
	if err = d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
