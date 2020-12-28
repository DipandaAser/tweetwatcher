package main

import (
	"encoding/json"
	"github.com/DipandaAser/tweetwatcher/bot"
	"github.com/DipandaAser/tweetwatcher/config"
	"github.com/DipandaAser/tweetwatcher/twitter"
	"io/ioutil"
	"log"
)

func main() {

	// Read project config file and creating a template if not exists
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		config.CreateConfigTemplate()
		log.Fatal("Missing configuration file. A blank Template will be create. Please fill it.")
	}

	var ProjectConfig config.Configuration

	err = json.Unmarshal(b, &ProjectConfig)
	if err != nil {
		log.Fatal("Configuration file is unreadable. A blank Template will be create. Please fill it.")
	}

	// Check requires fields
	err = ProjectConfig.Check()
	if err != nil {
		log.Fatal(err)
	}

	twitter.ProjectConfig = &ProjectConfig
	bot.ProjectConfig = &ProjectConfig

	// Start bot
	go bot.Start()

	// GetTweets start fetch all tweets
	twitter.GetTweets()
}
