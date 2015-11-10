package service

import (
	"bytes"
	"crypto/tls"
	"errors"
	"log"
	"net/smtp"
	"os"
	"text/template"
)

type SmtpTemplateData struct {
	From    string
	To      string
	Subject string
	Body    string
	Sender  string
}

const (
	MAIL_TEMPLATE = `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}

{{.Body}}

Sincerely,

{{.Sender}}
`

	MAIL_FROM              = "info@locals.ie"
	MAIL_HOST              = "server3.mywebdesign.ie"
	MAIL_PORT              = "25"
	MAIL_TEMPLATE_INTEREST = "Thanks for your interest in Locals.ie. \n\r We will be adding features overtime and will let you know about them as they become ready to use. \n\r"
)

type MailSender interface {
	Send(from, to, content string) error
}

type DefaultSender struct {
}

func (DefaultSender) Send(from, to, content string) error {
	var (
		err  error
		c    *smtp.Client
		auth smtp.Auth
	)
	log.Println("mail: sending mail to " + to)
	auth, err = DefaultSender{}.auth()
	if nil != err {
		return err
	}
	c, err = smtp.Dial(MAIL_HOST + ":" + MAIL_PORT)
	if nil != err {
		return err
	}
	defer c.Close()

	tlc := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         MAIL_HOST,
	}
	if err = c.StartTLS(tlc); err != nil {
		log.Println("tls error " + err.Error())
		return err
	}
	c.Auth(auth)
	c.Mail(from)
	c.Rcpt(to)

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}
	defer wc.Close()
	buf := bytes.NewBufferString(content)
	if _, err = buf.WriteTo(wc); err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Println("email: failed to write " + err.Error())
		return err
	}
	return err
}

func (DefaultSender) auth() (smtp.Auth, error) {
	var user string = os.Getenv("MAIL_USER")
	var pass string = os.Getenv("MAIL_PASS")
	var auth smtp.Auth

	if "" == user || "" == pass {
		return nil, errors.New("missing user or pass ensure env MAIL_USER and MAIL_PASS set")
	}
	auth = smtp.PlainAuth("", user, pass, MAIL_HOST)
	return auth, nil
}

func NewMailSender() MailSender {
	return DefaultSender{}
}

func parseTemplate(t, from, to, subject, sender string) (string, error) {

	var (
		err error
		doc bytes.Buffer
	)

	context := &SmtpTemplateData{
		From:    from,
		To:      to,
		Subject: subject,
		Body:    t,
		Sender:  "Craig",
	}
	parsed := template.New("emailTemplate")

	parsed, err = parsed.Parse(MAIL_TEMPLATE)
	if err != nil {
		log.Print("error trying to parse mail template")
	}
	err = parsed.Execute(&doc, context)
	return string(doc.Bytes()), err
}

func SendMailTemplate(template string, sender MailSender, from, to string) error {
	var (
		err     error
		content string
	)

	if MAIL_TEMPLATE_INTEREST == template {
		content, err = parseTemplate(MAIL_TEMPLATE_INTEREST, from, to, "Thanks for your interest in locals.ie", "Craig")
	} else {
		err = errors.New("no such template")
	}
	if nil != err {
		return err
	}
	return sender.Send(from, to, content)
}
