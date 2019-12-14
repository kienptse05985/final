package main

import (
	"fmt"
	"github.com/robfig/cron"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Container struct {
	MongoClient *mongo.Client
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
	clientOptions := options.ClientOptions{
		Auth: &options.Credential{
			Username:   config.MongoUser,
			Password:   config.MongoPassword,
			AuthSource: config.MongoSource,
		},
	}
	client, err := mongo.NewClient(clientOptions.ApplyURI(fmt.Sprintf("mongodb://%s", config.MongoServer)))
	if err != nil {
		return err
	}
	container.MongoClient = client
	container.CronDaemon = cron.New()
	container.MailRepository = &SendMailHandler{
		Username: config.MailUserName,
		Password: config.MailPassword,
	}
	return nil
}
