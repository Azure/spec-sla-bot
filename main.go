package main

import (
	"log"
	"playground/listpr/email"
)

func main() {
	/*result, err := github.PullRequests()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Number of Pull Requests: %d\n", len(result.Items))
	for _, item := range result.Items {
		fmt.Printf("#%-5d %9.9s Created: %.55s %.55s\n",
			item.Number, item.User.Login, item.CreatedAt, item.Title)
	}*/
	//Send a test email with all requests
	errEmail := email.SendEmailToAssignee()
	if errEmail != nil {
		log.Fatal(errEmail)
	}
	//template.GenerateTemplate()
}
