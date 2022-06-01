package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DipandaAser/tweetwatcher/bot/subscriber"
	"github.com/DipandaAser/tweetwatcher/config"
	tb "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
)

const (
	DefaultApiURL       = "https://api.telegram.org"
	channelStartCommand = "start watch"
	channelStopCommand  = "stop watch"
)

type Command struct {
	Text        string                    `json:"command"`
	Description string                    `json:"description"`
	Handler     func(message *tb.Message) `json:"-"`
}

func SetCommands() error {
	url := DefaultApiURL + "/bot" + myBot.Token + "/" + "setMyCommands"

	cmds := []Command{}
	for _, command := range getCommands() {
		cmds = append(cmds, command)
	}

	data, _ := json.Marshal(cmds)
	params := map[string]string{
		"commands": string(data),
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(params); err != nil {
		return err
	}

	_, err := http.DefaultClient.Post(url, "application/json", &buf)
	if err != nil {
		return err
	}

	return nil
}

func getCommands() map[string]Command {
	return map[string]Command{
		"start": {
			Text:        "start",
			Description: fmt.Sprintf("Activate the receptions of %s tweets.", config.ProjectConfig.Hashtag),
			Handler:     startCommandHAndler,
		},
		"stop": {
			Text:        "stop",
			Description: fmt.Sprintf("Desactivate the receptions of %s tweets.", config.ProjectConfig.Hashtag),
			Handler:     stopCommandHandler,
		},
		"help": {
			Text:        "help",
			Description: "Display all bot commands",
			Handler:     helpCommandHandler,
		},
	}
}

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
				message = fmt.Sprintf("Hey @%s ! \nYour status has change you will now receive tweets with #%s . \nBy @iamdipanda", userName, config.ProjectConfig.Hashtag)
			} else {
				message = fmt.Sprintf("Hey @%s ! \nAn error occur when we change your status. \nTry again later. \nBy @iamdipanda", userName)
			}
		} else {
			message = fmt.Sprintf("Hey @%s ! \nYour status has already set to receive tweets with #%s . \nBy @iamdipanda", userName, config.ProjectConfig.Hashtag)
		}
	} else {
		_, err := subscriber.New(userName, chatID)
		if err == nil {
			message = fmt.Sprintf("Hey @%s ! \nYou will now receive tweets with #%s . \nBy @iamdipanda", userName, config.ProjectConfig.Hashtag)
		} else {
			message = fmt.Sprintf("Hey @%s ! \nAn error occur. \nTry again later. \nBy @iamdipanda", userName)
		}
	}

	// We send the message
	msg := tb.NewMessage(m.Chat.ID, message)
	_, _ = myBot.Send(msg)
}

func stopCommandHandler(m *tb.Message) {

	var message string
	userName := getUserNameOrTitle(m)

	chatID := fmt.Sprintf("%d", m.Chat.ID)
	if subscriber.IsExist(chatID) {
		if subscriber.IsActivate(chatID) {
			if err := subscriber.DesactivateSubscriber(chatID); err == nil {
				message = fmt.Sprintf("Hey @%s ! \nYour status has change you will not receive tweets with #%s . \nBy @iamdipanda", userName, config.ProjectConfig.Hashtag)
			} else {
				message = fmt.Sprintf("Hey @%s ! \nAn error occur when we change your status. \nTry again later. \nBy @iamdipanda", userName)
			}
		} else {
			message = fmt.Sprintf("Hey @%s ! \nYour status has already set to not receive tweets with #%s . \nBy @iamdipanda", userName, config.ProjectConfig.Hashtag)
		}
	} else {
		message = fmt.Sprintf("Hey @%s ! \n You have not yet enable the reception of tweets with #%s . Send the /start command to enable it.  \n By @iamdipanda", userName, config.ProjectConfig.Hashtag)
	}

	// We send the message
	msg := tb.NewMessage(m.Chat.ID, message)
	_, _ = myBot.Send(msg)
}

func helpCommandHandler(m *tb.Message) {
	message := "Welcome to Bot Help \n \n"

	for _, command := range getCommands() {
		message += fmt.Sprintf("- /%s : %s \n", command.Text, command.Description)
	}

	// For channel help
	message += fmt.Sprintf(
		"For Channel Help\n"+
			"If the bot is in a channel\n"+
			"- write << %s >> to start\n"+
			"- write << %s >> to stop\n", channelStartCommand, channelStopCommand)

	message += "By @iamdipanda"

	// We send the message
	msg := tb.NewMessage(m.Chat.ID, message)
	_, _ = myBot.Send(msg)
}

func getUserNameOrTitle(m *tb.Message) string {
	if m.Chat.IsGroup() || m.Chat.IsChannel() {
		return m.Chat.Title
	}
	return m.Chat.UserName
}
