package email

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Azure/spec-sla-bot/template"
	gomail "gopkg.in/gomail.v2"
)

//SendEmailToAssignee sends an email to a list of users
func SendEmailToAssignee() error {
	template.GenerateTemplate()
	b, err := ioutil.ReadFile("finalTemplate.html")
	if err != nil {
		fmt.Print(err)
	}
	str := string(b) // convert content to a 'string'
	//fmt.Print(str)
	m := gomail.NewMessage()
	m.SetHeader("From", "t-jaelli@microsoft.com")
	m.SetHeader("To", "t-jaelli@microsoft.com")
	m.SetHeader("Subject", "TEST")
	m.SetBody("text/html", str)

	user := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	//Send the email
	d := gomail.NewDialer("smtp.office365.com", 587, user, password)
	if err = d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
