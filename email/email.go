package email

import (
	"bytes"
	"gopkg.in/gomail.v2"
	"html/template"
)

type Config struct {
	ServerHost string
	ServerPort string
	SenderAddr string
	Username string
	Password string
}

type EmailSender struct {
	conf *Config
	template string
}

type Sender interface {
	Send(to []string) error
	ParseTemplate(filepath string, data interface{}) error
}

func NewEmailSender(conf *Config, template string) Sender {
	return &EmailSender{conf,template }
}

func (e *EmailSender) ParseTemplate(filepath string, data interface{}) error {

	t, err := template.ParseFiles(filepath)

	 if err != nil {
	 	return err
	 }

	buf := new(bytes.Buffer)

	if err := t.Execute(buf, data); err != nil {
		return err
	}

	e.template = buf.String()

	return nil
}

func (e *EmailSender) Send(to []string) error {

	m := gomail.NewMessage()
	m.SetHeader("From", e.conf.SenderAddr)
	m.SetHeader("To", "danyal.iqbal@cheveo.de")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", e.template)

	d := gomail.NewDialer("smtp-mail.outlook.com", 587, e.conf.Username, e.conf.Password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func SetupEmailCredentials(host, port, senderAddr, username, password string) Config {
	return Config {
		host,port, senderAddr, username, password,
	}
}


