package config

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"strconv"
)

type Configuration struct {
	Hashtag    string
	BotToken   string
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

func Init() {
	ProjectConfig.Hashtag = os.Getenv("HASHTAG")
	ProjectConfig.BotToken = os.Getenv("BOT_TOKEN")
	ProjectConfig.DBName = os.Getenv("DB_NAME")
	ProjectConfig.MongodbURI = os.Getenv("MONGO_URI")
	delay, err := strconv.Atoi(os.Getenv("SCRAP_DELAY"))
	if err == nil {
		ProjectConfig.ScrapDelay = delay
	}

	err = ProjectConfig.check()
	if err != nil {
		log.Fatal(err)
	}
}

// Check check if requires fields are provided
func (config *Configuration) check() error {
	if config.Hashtag == "" {
		return errors.New("please specify a \"HASHTAG\" env var")
	}
	if config.BotToken == "" {
		return errors.New("please specify a \"BOT_TOKEN\" env var")
	}
	if config.MongodbURI == "" {
		return errors.New("please specify a \"MONGO_URI\" env var")
	}
	if config.DBName == "" {
		return errors.New("please specify the \"DB_NAME\" env var")
	}
	return nil
}
