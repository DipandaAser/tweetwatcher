package bot

import (
	"fmt"
	"github.com/DipandaAser/tweetwatcher/bot/subscriber"
	"github.com/DipandaAser/tweetwatcher/config"
	tb "gopkg.in/tucnak/telebot.v2"
)

func startCommandHAndler(m *tb.Message) {
	var message string
	userName := getUserNameOrTitle(m)

	//We try to fetch if that chat_id is already in the subscriber list and if the status of receiving tweet is ok or not
	//If not we add it with the status ok
	//If yes we say him(her) he already subscribe to receive tweet

	chatID := fmt.Sprintf("%d", m.Chat.ID)

	if subscriber.IsExist(chatID) {
		if !subscriber.IsActivate(chatID) {
			if err := subscriber.ActivateSubscriber(chatID); err == nil {
				message = fmt.Sprintf("Hey @%s ! \nYour status has change you will now receive all %s tweets. \nBy @iamdipanda", userName, config.ProjectConfig.Hashtag)
			} else {
				message = fmt.Sprintf("Hey @%s ! \nAn error occur when we change your status. \nTry again later. \nBy @iamdipanda", userName)
			}
		} else {
			message = fmt.Sprintf("Hey @%s ! \nYour status has already set to receive all %s tweets. \nBy @iamdipanda", userName, config.ProjectConfig.Hashtag)
		}
	} else {
		_, err := subscriber.New(userName, chatID)
		if err == nil {
			message = fmt.Sprintf("Hey @%s ! \nYou will now receive all %s tweets. \nBy @iamdipanda", userName, config.ProjectConfig.Hashtag)
		} else {
			message = fmt.Sprintf("Hey @%s ! \nAn error occur. \nTry again later. \nBy @iamdipanda", userName)
		}
	}

	// We send the message
	_, _ = myBot.Send(m.Chat, message)
}

func stopCommandHandler(m *tb.Message) {

	var message string
	userName := getUserNameOrTitle(m)

	chatID := fmt.Sprintf("%d", m.Chat.ID)
	if subscriber.IsExist(chatID) {
		if subscriber.IsActivate(chatID) {
			if err := subscriber.DesactivateSubscriber(chatID); err == nil {
				message = fmt.Sprintf("Hey @%s ! \nYour status has change you will not receive all %s tweets. \nBy @iamdipanda", userName, config.ProjectConfig.Hashtag)
			} else {
				message = fmt.Sprintf("Hey @%s ! \nAn error occur when we change your status. \nTry again later. \nBy @iamdipanda", userName)
			}
		} else {
			message = fmt.Sprintf("Hey @%s ! \nYour status has already set to not receive all %s tweets. \nBy @iamdipanda", userName, config.ProjectConfig.Hashtag)
		}
	} else {
		message = fmt.Sprintf("Hey @%s ! \n You have not yet enable the reception of all %s tweets. Send the /start command to enable it.  \n By @iamdipanda", userName, config.ProjectConfig.Hashtag)
	}

	// We send the message
	_, _ = myBot.Send(m.Chat, message)
}

func helpCommandHandler(m *tb.Message) {
	message := "Welcome to Bot Help \n \n"

	for _, command := range commands {
		message += fmt.Sprintf("- %s : %s \n", command.Text, command.Description)
	}

	// For channel help
	message += fmt.Sprintf(
		"For Channel Help\n"+
			"If the bot is in a channel\n"+
			"- write << %s >> to start\n"+
			"- write << %s >> to stop\n", channelStartCommand, channelStopCommand)

	message += "By @iamdipanda"

	// We send the message
	_, _ = myBot.Send(m.Chat, message)
}

func getUserNameOrTitle(m *tb.Message) string {
	if m.FromGroup() || m.FromChannel() {
		return m.Chat.Title
	}
	return m.Sender.Username
}
