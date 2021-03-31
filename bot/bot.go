package bot

import (
	"fmt"
	"github.com/DipandaAser/tweetwatcher/bot/subscriber"
	"github.com/DipandaAser/tweetwatcher/config"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
)

var myBot *tb.Bot
var commands []tb.Command

func Start() {
	webhook := &tb.Webhook{
		Listen:   ":" + config.ProjectConfig.Port,
		Endpoint: &tb.WebhookEndpoint{PublicURL: config.ProjectConfig.PublicURL},
	}

	b, err := tb.NewBot(tb.Settings{
		Token:   config.ProjectConfig.BotToken,
		Poller:  webhook,
		Verbose: false,
	})

	if err != nil {
		log.Fatal(err)
	}

	myBot = b

	// Setting all bot commands and create handler for each command
	commands = []tb.Command{
		{Text: "/start", Description: fmt.Sprintf("Activate the receptions of %s tweets.", config.ProjectConfig.Hashtag)},
		{Text: "/stop", Description: fmt.Sprintf("Desactivate the receptions of %s tweets.", config.ProjectConfig.Hashtag)},
		{Text: "/help", Description: "Display all bot commands"},
	}
	_ = myBot.SetCommands(commands)

	myBot.Handle("/start", startCommandHAndler)

	myBot.Handle("/stop", stopCommandHandler)

	myBot.Handle("/help", helpCommandHandler)

	log.Println("#### Bot start sucessfully ####")
	myBot.Start()
}

func BulkSendPhoto(tweetScreenshotUrl string, tweetLink string, tweetUserName string) {

	caption := fmt.Sprintf("New Tweet from : @%s .\n\n-------\n\nLink: %s ", tweetUserName, tweetLink)
	photoMsg := &tb.Photo{Caption: caption, File: tb.FromURL(tweetScreenshotUrl)}

	chats, err := subscriber.GetActivatedSubscribers()
	if err != nil {
		return
	}

	for _, chat := range chats {
		_, _ = myBot.Send(chat, photoMsg)
	}

}

func BulkSendText(description string, tweetLink string, tweetUserName string) {

	msg := fmt.Sprintf("New Tweet from : %s .\n\n-------\n\n %s \n\n-------\n\nLink: %s ", tweetUserName, description, tweetLink)
	chats, err := subscriber.GetActivatedSubscribers()
	if err != nil {
		return
	}

	for _, chat := range chats {
		_, _ = myBot.Send(chat, msg)
	}

}
