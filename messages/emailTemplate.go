package messages

import (
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/Azure/spec-sla-bot/models"
)

//CreatePrimaryTemplate populates the finalPrimaryTemplate for the assignee email
func CreatePrimaryTemplate(info *MessageContent) error {
	fmap := template.FuncMap{
		"FormatNumber":   FormatNumber,
		"FormatAssignee": FormatAssignee}
	t := template.Must(template.New("primaryTemplate.tmpl").Funcs(fmap).ParseFiles("./templates/primaryTemplate.tmpl"))
	handle, err := os.Create("finalPrimaryTemplate.html")
	err = t.Execute(handle, info)
	if err != nil {
		return err
	}
	return nil
}

//CreateManagerTemplate should this accept something other than messageContent
func CreateManagerTemplate(emails []models.Email) error {
	fmap := template.FuncMap{
		"FormatNumber":   FormatNumber,
		"FormatAssignee": FormatAssignee}
	t := template.Must(template.New("managerTemplate.tmpl").Funcs(fmap).ParseFiles("./templates/managerTemplate.tmpl"))
	handle, err := os.Create("finalManagerTemplate.html")
	err = t.Execute(handle, emails)
	if err != nil {
		return err
	}
	return nil
}

/*func CreateAssigneeTemplate() {
	//Map of function names to functions
	//box := packr.NewBox("../templates")
	fmap := template.FuncMap{
		"FormatNumber":   FormatNumber,
		"FormatUser":     FormatUser,
		"FormatAssignee": FormatAssignee,
		"FormatTime":     FormatTime,
		"FormatTitle":    FormatTitle}
	t := template.Must(template.New("assigneeTemplate.tmpl").Funcs(fmap).ParseFiles("./templates/assigneeTemplate.tmpl"))
	result, err := PullRequests()
	if err != nil {
		log.Fatal(err)
	}
	handle, err := os.Create("finalTemplate.html")
	err = t.Execute(handle, *result)
	if err != nil {
		panic(err)
	}
}*/

func FormatNumber(number int) string {
	return fmt.Sprintf("#%-5d", number)
}

func FormatUser(user User) string {
	return fmt.Sprintf("%9.9s", user.Login)
}

func FormatAssignee(assignee Assignee) string {
	return fmt.Sprintf("%9.9s", assignee.Login)
}

func FormatTime(createdAt time.Time) string {
	return fmt.Sprintf("%.55s", createdAt)
}

func FormatTitle(title string) string {
	return fmt.Sprintf("%.55s\n", title)
}
