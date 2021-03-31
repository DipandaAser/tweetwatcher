package bot

import (
	"fmt"
	"github.com/DipandaAser/tweetwatcher/bot/subscriber"
	"github.com/DipandaAser/tweetwatcher/config"
	tb "gopkg.in/tucnak/telebot.v2"
)

func startCommandHAndler(m *tb.Message) {
	var message string

	//We try to fetch if that chat_id is already in the suscriber list and if the status of receiving tweet is ok or not
	//If not we add it with the status ok
	//If yes we say him(her) he already suscribe to receive tweet

	chatID := fmt.Sprintf("%d", m.Chat.ID)
	if subscriber.IsExist(chatID) {
		if !subscriber.IsActivate(chatID) {
			if err := subscriber.ActivateSubscriber(chatID); err == nil {
				message = fmt.Sprintf("Hey @%s ! \nYour status has change you will now receive all %s tweets. \nBy @iamdipanda", m.Chat.Username, config.ProjectConfig.Hashtag)
			} else {
				message = fmt.Sprintf("Hey @%s ! \nAn error occur when we change your status. \nTry again later. \nBy @iamdipanda", m.Chat.Username)
			}
		} else {
			message = fmt.Sprintf("Hey @%s ! \nYour status has already set to receive all %s tweets. \nBy @iamdipanda", m.Chat.Username, config.ProjectConfig.Hashtag)
		}
	} else {
		_, err := subscriber.New(m.Sender.Username, chatID)
		if err == nil {
			message = fmt.Sprintf("Hey @%s ! \nYou will now receive all %s tweets. \nBy @iamdipanda", m.Chat.Username, config.ProjectConfig.Hashtag)
		} else {
			message = fmt.Sprintf("Hey @%s ! \nAn error occur. \nTry again later. \nBy @iamdipanda", m.Chat.Username)
		}
	}

	// We send the message
	_, _ = myBot.Send(m.Sender, message)
}

func stopCommandHandler(m *tb.Message) {

	var message string

	chatID := fmt.Sprintf("%d", m.Chat.ID)
	if subscriber.IsExist(chatID) {
		if subscriber.IsActivate(chatID) {
			if err := subscriber.DesactivateSubscriber(chatID); err == nil {
				message = fmt.Sprintf("Hey @%s ! \nYour status has change you will not receive all %s tweets. \nBy @iamdipanda", m.Chat.Username, config.ProjectConfig.Hashtag)
			} else {
				message = fmt.Sprintf("Hey @%s ! \nAn error occur when we change your status. \nTry again later. \nBy @iamdipanda", m.Chat.Username)
			}
		} else {
			message = fmt.Sprintf("Hey @%s ! \nYour status has already set to not receive all %s tweets. \nBy @iamdipanda", m.Chat.Username, config.ProjectConfig.Hashtag)
		}
	} else {
		message = fmt.Sprintf("Hey @%s ! \n You have not yet enable the reception of all %s tweets. Send the /start command to enable it.  \n By @iamdipanda", m.Chat.Username, config.ProjectConfig.Hashtag)
	}

	// We send the message
	_, _ = myBot.Send(m.Chat, message)
}

func helpCommandHandler(m *tb.Message) {
	message := "Welcome to Bot Help \n \n"

	for _, command := range commands {
		message += fmt.Sprintf("- %s : %s \n", command.Text, command.Description)
	}

	message += "By @iamdipanda"

	// We send the message
	_, _ = myBot.Send(m.Chat, message)
}
