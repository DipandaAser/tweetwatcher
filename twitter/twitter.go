package twitter

import (
	"context"
	"errors"
	"fmt"
	"github.com/DipandaAser/tweetwatcher/bot"
	"github.com/DipandaAser/tweetwatcher/config"
	twitterScraper "github.com/n0madic/twitter-scraper"
	"log"
	"net/http"
	"time"
)

const minScrappingDelayInSecond = 10

// GetTweets will fetch all tweets with the hashtag define in config file
func GetTweets() {

	log.Println("#### Start fetching tweet ####")
	scrapper := twitterScraper.New()
	scrapper.SetSearchMode(twitterScraper.SearchLatest)

	for tweetResult := range scrapper.SearchTweets(context.TODO(), fmt.Sprintf("#%s", config.ProjectConfig.Hashtag), 100) {

		if tweetResult.Error != nil {
			continue
		}

		// If the tweet already exist we don't save and publish it
		if IsExist(tweetResult.ID) {
			continue
		}

		log.Println("New tweet found")

		date, _ := tweetResult.TimeParsed.MarshalText()
		_, _ = Save(tweetResult.ID, tweetResult.Text, tweetResult.Username, tweetResult.PermanentURL, string(date))

		// If error occur when taking screenshot, we send text message instead of photo message
		photo, err := takeScreenshot(tweetResult.ID)
		if err == nil {
			go bot.BulkSendPhoto(photo, tweetResult.PermanentURL, tweetResult.Username)
		} else {
			go bot.BulkSendText(tweetResult.Text, tweetResult.PermanentURL, tweetResult.Username)
		}
	}

	log.Println("#### Sleep ####")
}

func GetScrapDelay() time.Duration {
	var scrapDelay time.Duration
	scrapDelay = minScrappingDelayInSecond * time.Second
	if config.ProjectConfig.ScrapDelay < minScrappingDelayInSecond {
		delay, err := time.ParseDuration(fmt.Sprintf("%vs", config.ProjectConfig.ScrapDelay))
		if err == nil {
			scrapDelay = delay
		}
	}

	return scrapDelay
}

func takeScreenshot(tweetId string) (string, error) {

	url := fmt.Sprintf("https://tweet2image.vercel.app/%s.png?lang=en&tz=0&theme=light&scale=1", tweetId)
	rep, err := http.Get(url)
	if err != nil {
		return "", err
	}

	if rep.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("request fail with status code %d", rep.StatusCode))
	}

	return url, nil
}
