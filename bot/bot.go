package bot

import (
	"fmt"
	"github.com/DipandaAser/tweetwatcher/bot/subscriber"
	"github.com/DipandaAser/tweetwatcher/config"
	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

var myBot *tb.BotAPI

func Start() {

	bot, err := tb.NewBotAPI(config.ProjectConfig.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	myBot = bot

	// We remove webhook to avoid conflict. Webhook can't exist when we use socket to connect our bot
	_, err = bot.RemoveWebhook()
	if err != nil {
		log.Fatal(err)
	}

	// Setting all bot commands
	_ = SetCommands()

	// Open the websocket and begin listening.
	u := tb.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("#### Bot start successfully ####")

	commands := getCommands()
	for update := range updates {

		// ignore any non-Message Updates
		if update.Message == nil {
			continue
		}

		// Check if the command is a registered command
		if update.Message.IsCommand() {
			cmd, cmdExist := commands[update.Message.Command()]
			if cmdExist {

				if cmd.Handler != nil {
					cmd.Handler(update.Message)
				}
			}
		}

		// For channel we use text to detect
		if update.Message.Chat.IsChannel() {
			switch update.Message.Text {
			case channelStartCommand:
				startCommandHAndler(update.Message)
				break
			case channelStopCommand:
				stopCommandHandler(update.Message)
				break
			default:
				break
			}
		}

	}

	log.Println("#### Bot shutdown successfully ####")
}

func BulkSendPhoto(tweetScreenshotUrl string, tweetScreenshot []byte, tweetLink string, tweetUserName string) {

	caption := fmt.Sprintf("New Tweet from : @%s .\n\n-------\n\nLink: %s ", tweetUserName, tweetLink)

	chats, err := subscriber.GetActivatedSubscribers()
	if err != nil {
		return
	}

	photo := tb.FileBytes{Bytes: tweetScreenshot}
	for _, chat := range chats {
		photoMsg := tb.NewPhotoUpload(chat.Recipient(), photo)
		photoMsg.Caption = caption
		_, _ = myBot.Send(photoMsg)
	}

}

func BulkSendText(description string, tweetLink string, tweetUserName string) {

	message := fmt.Sprintf("New Tweet from : @%s .\n\n-------\n\n %s \n\n-------\n\nLink: %s ", tweetUserName, description, tweetLink)
	chats, err := subscriber.GetActivatedSubscribers()
	if err != nil {
		return
	}

	for _, chat := range chats {
		msg := tb.NewMessage(chat.Recipient(), message)
		_, _ = myBot.Send(msg)
	}

}
