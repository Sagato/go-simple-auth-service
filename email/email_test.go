package email

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestEmail_ParseTemplate(t *testing.T) {

	emailConfig := SetupEmailCredentials(
		os.Getenv("serverHost"),
		os.Getenv("serverPort"),
		os.Getenv("senderAddress"),
		os.Getenv("username"),
		os.Getenv("password"))

	fmt.Println(emailConfig)

	es := &EmailSender {
		&emailConfig,
		"Test",
	}

	wd, err := os.Getwd()

	// Get working Directory for correct html template file paths
	if err != nil {
		t.Error(err.Error())
	}

	// Init Anonymous struct with template data
	data := struct {
		Name string
		Greetings string
		WebsiteUrl string
		Email string
		WebsiteName string
		Year string
		CompanyName string
		ActivationLink string
		From string
	}{
		"sagat",
		"Hello",
		"https://www.example.com",
		"sagat@web.de",
		"Sagateyson",
		"2019",
		"sagat corp.",
		"https://www.example.com/activate",
		"testdomain.com",
	}

	if err := es.ParseTemplate(filepath.Join(wd,  "../email/templates/registration_mail.html"), data); err != nil {
		t.Error(err.Error())
	}

	if es.template == "" {
		t.Error("template wasnt parsed!")
	}
}
