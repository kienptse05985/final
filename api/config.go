package main

type Config struct {
	LogLevel int

	InternalAPI string
	Binding     string

	GoogleReCaptchaSecret string

	MailUserName string
	MailPassword string

}
