package messages

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/Azure/spec-sla-bot/models"
	gomail "gopkg.in/gomail.v2"
)

//Email information structure
type Email struct {
	Password     string
	EmailAddress string
	Port         int
}

//SendEmailToAssignee sends an email to a list of users
func SendEmailToAssignee(ctx context.Context, info *MessageContent) error {
	err := CreatePrimaryTemplate(info)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile("finalPrimaryTemplate.html")
	if err != nil {
		fmt.Print(err)
		return err
	}

	//Determine email recipient from database
	queryString := fmt.Sprintf("SELECT EmailLogin FROM [User] WHERE GitHubUser = '%s';", info.AssigneeLogin)
	rows, err := InfraDB.QueryContext(ctx, queryString)
	if err != nil {
		return err
	}
	defer rows.Close()
	var emailTo string
	if rows != nil {
		for rows.Next() {
			err = rows.Scan(&emailTo)
			if err != nil {
				return err
			}
		}
		//Can be removed once we use a default email in the repo or use the infra monitor database
		//to send the email
		if emailTo == "" {
			emailTo = "t-jaelli@microsoft.com"
		}
	} else {
		return fmt.Errorf("Cannot find email for %s to send SLA reminder email", info.AssigneeLogin)
	}

	//Parse email connection string
	emailInfo, err := parseEmailURL()
	if err != nil {
		return err
	}

	//Format email
	str := string(b)
	m := gomail.NewMessage()
	m.SetHeader("From", emailInfo.EmailAddress)
	m.SetHeader("To", emailTo)
	m.SetHeader("Subject", "SLA Violation - Outstanding Pull Request")
	m.SetBody("text/html", str)

	d := gomail.NewDialer("smtp.office365.com", emailInfo.Port, emailInfo.EmailAddress, emailInfo.Password)
	if err = d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func SendEmailToManager(ctx context.Context, info *MessageContent) error {
	//Get record of emails sent during date ran
	emails := []models.Email{}
	err := models.DB.RawQuery(`SELECT * FROM emails WHERE [time_sent] > '?' AND [time_sent] <= '?'`, time.Now().AddDate(0, 0, -7), time.Now()).All(&emails)
	if err != nil {
		log.Print("Could not make query")
		return err
	}
	err = CreateManagerTemplate(emails)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile("finalManagerTemplate.html")
	if err != nil {
		fmt.Print(err)
		return err
	}

	//Parse email connection string
	emailInfo, err := parseEmailURL()
	if err != nil {
		return err
	}

	str := string(b)
	m := gomail.NewMessage()
	m.SetHeader("From", emailInfo.EmailAddress)
	//Will always be sent to the manager
	m.SetHeader("To", "t-jaelli@microsoft.com")
	m.SetHeader("Subject", "SLA Violations - A Week In Review")
	m.SetBody("text/html", str)

	d := gomail.NewDialer("smtp.office365.com", emailInfo.Port, emailInfo.EmailAddress, emailInfo.Password)
	if err = d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func parseEmailURL() (*Email, error) {
	emailURL := os.Getenv("CUSTOMCONNSTR_EMAIL_URL")

	parsed, err := url.Parse(emailURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	port, _ := strconv.Atoi(parsed.Port())
	password, _ := parsed.User.Password()
	emailStruct := &Email{
		Password:     password,
		Port:         port,
		EmailAddress: "t-jaelli@microsoft.com",
	}
	return emailStruct, nil
}
