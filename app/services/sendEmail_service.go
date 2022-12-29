package services

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to string, subject, body string) error {
	// Connect to the Gmail SMTP server
	from := os.Getenv("GMAIL_ID")
	password := os.Getenv("GMAIL_PASSWORD")
	host := "smtp.gmail.com"
	port := "587"
	subject = "Subject: " + subject + "\n"
	address := host + ":" + port
	message := []byte(subject + body)
	auth := smtp.PlainAuth("", from, password, host)
	err := smtp.SendMail(address, auth, from, []string{to}, message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Email Sent!")
	return nil
}
