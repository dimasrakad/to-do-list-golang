package utils

import (
	"bytes"
	"html/template"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"to-do-list-golang/config"
	"to-do-list-golang/models"

	"gopkg.in/gomail.v2"
)

var (
	emailTemplate     *template.Template
	loadTemplateOnce  sync.Once
	emailTemplatePath string
)

func SendNotificationEmail(notificationType string, todo models.Todo) error {
	var err error
	cfg := config.LoadConfig()

	subject := "Reminder: To Do \"" + todo.Title + "\" is due in " + notificationType

	loadTemplateOnce.Do(func() {
		emailTemplatePath = filepath.Join("templates", "emails", "reminder.html")
		emailTemplate, err = template.ParseFiles(emailTemplatePath)
		if err != nil {
			log.Println("Error parsing email template:", err)
		}
	})
	if err != nil {
		return err
	}

	smtpPort, err := strconv.Atoi(cfg.SMTPPort)
	if err != nil {
		log.Println("Error parsing SMTP Port:", err)
		return err
	}

	d := gomail.NewDialer(cfg.SMTPHost, smtpPort, cfg.SMTPUser, cfg.SMTPPass)

	s, err := d.Dial()
	if err != nil {
		log.Println("Error dialing SMTP:", err)
		return err
	}
	defer s.Close()

	var wg sync.WaitGroup
	for _, assignee := range todo.Assignees {
		wg.Add(1)
		go func(assignee models.User) {
			defer wg.Done()

			data := map[string]string{
				"Name":             assignee.Name,
				"Title":            todo.Title,
				"Due":              todo.Due.Format("02 Jan 2006 15:04"),
				"NotificationType": notificationType,
			}

			var body bytes.Buffer
			if err := emailTemplate.Execute(&body, data); err != nil {
				log.Printf("Error executing email template for %s: %v\n", assignee.Email, err)
				return
			}

			mail := gomail.NewMessage()
			mail.SetHeader("From", cfg.SMTPFrom)
			mail.SetHeader("To", assignee.Email)
			mail.SetHeader("Subject", subject)
			mail.SetBody("text/html", body.String())

			if err := gomail.Send(s, mail); err != nil && !strings.Contains(err.Error(), "250 2.0.0 OK") {
				log.Printf("Error sending email to %s: %v\n", assignee.Email, err)
			}
		}(assignee)
	}

	wg.Wait()
	return nil
}
