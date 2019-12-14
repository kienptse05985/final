package main

import (
	"github.com/robfig/cron"
)

type Container struct {
	Config      Config

	CronDaemon     *cron.Cron
	MailRepository *SendMailHandler
}

func NewContainer(config Config) (*Container, error) {
	container := new(Container)
	err := container.InitContainer(config)
	if err != nil {
		return nil, err
	}
	return container, nil
}

func (container *Container) InitContainer(config Config) (err error) {

	// load dependencies
	if err := container.LoadDependencies(config); err != nil {
		return err
	}

	return
}

func (container *Container) LoadDependencies(config Config) (err error) {
	container.Config = config
	container.CronDaemon = cron.New()
	container.MailRepository = &SendMailHandler{
		Username: config.MailUserName,
		Password: config.MailPassword,
	}
	return nil
}
