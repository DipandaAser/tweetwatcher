package main

import (
	"context"
	"github.com/DipandaAser/tweetwatcher/bot"
	"github.com/DipandaAser/tweetwatcher/config"
	"github.com/DipandaAser/tweetwatcher/twitter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {

	config.ProjectConfig.Hashtag = os.Getenv("HASHTAG")
	config.ProjectConfig.BotToken = os.Getenv("BOT_TOKEN")
	config.ProjectConfig.Port = os.Getenv("PORT")
	config.ProjectConfig.DBName = os.Getenv("DB_NAME")
	config.ProjectConfig.MongodbURI = os.Getenv("MONGO_URI")
	config.ProjectConfig.PublicURL = os.Getenv("PUBLIC_URL")
	delay, err := strconv.Atoi(os.Getenv("SCRAP_DELAY"))
	if err == nil {
		config.ProjectConfig.ScrapDelay = delay
	}

	ctx := context.TODO()
	config.MongoCtx = &ctx

	// Check requires var to start program
	err = config.ProjectConfig.Check()
	if err != nil {
		log.Fatal(err)
	}

	// ─── MONGO ──────────────────────────────────────────────────────────────────────
	err = MongoConnect()
	if err != nil {
		log.Fatal("Can't setup mongodb")
	}

	// ─── WE REFRESH THE MONGO CONNECTION EACH 10MINS ──────────────────────────────────────
	mongoTicker := time.NewTicker(time.Minute * 10)
	defer mongoTicker.Stop()
	go func() {
		for range mongoTicker.C {
			go MongoReconnectCheck()
		}
	}()

	// Start fetch tweets
	scrapperTicker := time.NewTicker(twitter.GetScrapDelay())
	defer scrapperTicker.Stop()
	go func() {
		for range scrapperTicker.C {
			go twitter.GetTweets()
		}
	}()

	// Start bot
	bot.Start()
}

// MongoConnect connects to mongoDB
func MongoConnect() error {

	clientOptions := options.Client().ApplyURI(config.ProjectConfig.MongodbURI)

	client, err := mongo.Connect(*config.MongoCtx, clientOptions)
	if err != nil {
		return err
	}

	// We make sure we have been connected
	err = client.Ping(*config.MongoCtx, readpref.Primary())
	if err != nil {
		return err
	}

	db := client.Database(config.ProjectConfig.DBName)
	config.DB = db

	return nil
}

// MongoReconnectCheck reconnects to MongoDB
func MongoReconnectCheck() {

	// We make sure we are still connected
	err := config.DB.Client().Ping(*config.MongoCtx, readpref.Primary())
	if err != nil {
		// We reconnect
		_ = MongoConnect()
	}
}
