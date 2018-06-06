package template

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/Azure/spec-sla-bot/github"
)

func GenerateTemplate() {
	//Map of function names to functions
	fmap := template.FuncMap{
		"FormatNumber":   FormatNumber,
		"FormatUser":     FormatUser,
		"FormatAssignee": FormatAssignee,
		"FormatTime":     FormatTime,
		"FormatTitle":    FormatTitle}
	t := template.Must(template.New("assigneeTemplate.tmpl").Funcs(fmap).ParseFiles("assigneeTemplate.tmpl"))
	result, err := github.PullRequests()
	if err != nil {
		log.Fatal(err)
	}
	handle, err := os.Create("finalTemplate.html")
	//fred, err := ioutil.TempFile()
	err = t.Execute(handle, *result)
	if err != nil {
		panic(err)
	}
}

/*func FormatPullRequest(item github.Request) string {
	formattedString := fmt.Sprintf("#%-5d %9.9s Created: %.55s %.55s\n",
		item.Number, item.User.Login, item.CreatedAt, item.Title)
	return formattedString
}*/

func FormatNumber(number int) string {
	return fmt.Sprintf("#%-5d", number)
}

func FormatUser(user github.User) string {
	return fmt.Sprintf("%9.9s", user.Login)
}

func FormatAssignee(assignee github.Assignee) string {
	return fmt.Sprintf("%9.9s", assignee.Login)
}

func FormatTime(createdAt time.Time) string {
	return fmt.Sprintf("%.55s", createdAt)
}

func FormatTitle(title string) string {
	return fmt.Sprintf("%.55s\n", title)
}
