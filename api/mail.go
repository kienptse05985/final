package main

import (
	"encoding/json"
	"fmt"
	"net/smtp"
)

type MailData struct {
	From        string   `json:"from"`
	To          []string `json:"to"`
	CC          []string `json:"cc"`
	BCC         []string `json:"bcc"`
	Subject     string   `json:"subject"`
	ContentType string   `json:"contentType"`
	Content     string   `json:"content"`
}

type SendMailHandler struct {
	Username string
	Password string
}

func (h *SendMailHandler) SendByGmail(body []byte) error {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	var data MailData
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("SendMailHandler: decoding body failed")
		return err
	}

	from := h.Username
	pass := h.Password

	to := data.To[0]
	msg := fmt.Sprintf("From: %s\nTo:%s\nSubject:%s\n%s\n\n%s", from, to, data.Subject, mime, data.Content)

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, data.To, []byte(msg))

	if err != nil {
		fmt.Println("SendMailHandler: sending mail failed")
		return err
	}

	fmt.Printf("SendMailHandler: successfully sent mail to %s\n", to)
	return nil
}

