package config

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type Configuration struct {
	Hashtag    string
	BotToken   string
	Port       string
	PublicURL  string
	MongodbURI string
	DBName     string
	// ScrapDelay represent the interval before the scrapper perform another scraping session
	ScrapDelay int
}

// ProjectConfig holds info of project settings
var ProjectConfig = &Configuration{}

// MongoCtx is the mongo context
var MongoCtx *context.Context

// DB is the mongo db
var DB *mongo.Database

// Check check if requires fields are provided
func (config *Configuration) Check() error {
	if config.Hashtag == "" {
		return errors.New("please specify a \"HASHTAG\" env var")
	}
	if config.BotToken == "" {
		return errors.New("please specify a \"BOT_TOKEN\" env var")
	}
	if config.Port == "" {
		return errors.New("please specify a \"PORT\" env var")
	}
	if config.PublicURL == "" {
		return errors.New("please specify the \"PUBLIC_URL\"  env var")
	}
	if config.MongodbURI == "" {
		return errors.New("please specify a \"MONGO_URI\" env var")
	}
	if config.DBName == "" {
		return errors.New("please specify the \"DB_NAME\" env var")
	}
	return nil
}
