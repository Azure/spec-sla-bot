package messages

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"

	gomail "gopkg.in/gomail.v2"
)

//SendEmailToAssignee sends an email to a list of users
func SendEmailToAssignee(ctx context.Context, info *Message) error {
	CreatePrimaryTemplate(info)
	b, err := ioutil.ReadFile("finalPrimaryTemplate.html")
	if err != nil {
		fmt.Print(err)
		return err
	}
	str := string(b)
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
	password, _ := parsed.User.Password()

	queryString := fmt.Sprintf("SELECT EmailLogin FROM [User] WHERE GitHubUser = '%s';", info.Assignee)
	fmt.Println("email selection query: ", queryString)
	rows, err := InfraDB.QueryContext(ctx, queryString)
	if err != nil {
		return err
	}
	defer rows.Close()
	var emailTo string
	if rows != nil {
		log.Print("rows does not equal null")
		for rows.Next() {
			err = rows.Scan(&emailTo)
			if err != nil {
				return err
			}
		}
		if emailTo == "" {
			emailTo = "t-jaelli@microsoft.com"
		}
		log.Printf("ISSUE")
		m.SetHeader("To", emailTo)
	} else {
		log.Printf("Cannot find email for %s to send SLA reminder email", info.Assignee)
		return err
	}

	m.SetHeader("From", "t-jaelli@microsoft.com")
	m.SetHeader("To", emailTo)
	m.SetHeader("Subject", "TEST")
	m.SetBody("text/html", str)
	d := gomail.NewDialer("smtp.office365.com", port, "t-jaelli@microsoft.com", password)
	if err = d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
