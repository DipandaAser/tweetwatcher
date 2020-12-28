package twitter

import (
	"context"
	"github.com/DipandaAser/tweetwatcher/bot"
	"github.com/DipandaAser/tweetwatcher/config"
	"github.com/DipandaAser/tweetwatcher/utils"
	"github.com/go-rod/rod"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"log"
)

var ProjectConfig *config.Configuration

// GetTweets will fetch all tweets with the hashtag define in config file
func GetTweets() {

	log.Println("#### Start getting tweet ####")
	twitterscraper.SetSearchMode(twitterscraper.SearchLatest)

	for {
		for tweet := range twitterscraper.SearchTweets(context.Background(), ProjectConfig.Hashtag, 50) {
			if tweet.Error != nil {
				continue
			}
			screenshotFile := tweet.ID + ".png"

			// If the tweet already exist we stop our current list of tweet and we restard search
			if utils.FileExists(screenshotFile) {
				break
			}

			log.Println("#### New tweet found ####")
			log.Printf("ID: %s\n", tweet.ID)
			log.Printf("User: %s\n", tweet.Username)
			log.Printf("Link: %s\n\n", tweet.PermanentURL)

			// If error occur when taking screenshoot, we send text message instead of photo message
			err := takeScreenshoot(tweet.PermanentURL, tweet.ID+".png")
			if err == nil {
				go bot.BulkSendPhoto(screenshotFile, tweet.PermanentURL, tweet.Username)
			} else {
				go bot.BulkSendText(tweet.Text, tweet.PermanentURL, tweet.Username)
			}
		}
	}
}

func takeScreenshoot(url string, filename string) error {
	browser := rod.New().MustConnect()
	defer browser.MustClose()
	page := browser.MustPage(url).MustWaitLoad()

	// We identified the tweet with this selector
	tweetNode, err := page.Element("article:last-of-type.css-1dbjc4n")
	if err != nil {
		return err
	}

	tweetNode.MustScreenshot(filename)
	return nil
}
