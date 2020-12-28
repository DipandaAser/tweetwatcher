package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DipandaAser/tweetwatcher/bot/suscriber"
	"github.com/DipandaAser/tweetwatcher/config"
	tb "gopkg.in/tucnak/telebot.v2"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

var ProjectConfig *config.Configuration

var myBot tb.Bot

func Start() {
	webhook := &tb.Webhook{
		Listen:   ":" + ProjectConfig.Port,
		Endpoint: &tb.WebhookEndpoint{PublicURL: ProjectConfig.PublicURL},
	}

	myBot, err := tb.NewBot(tb.Settings{
		Token:   ProjectConfig.BotToken,
		Poller:  webhook,
		Verbose: false,
	})

	if err != nil {
		log.Fatal(err)
	}

	// Setting all bot commands and create handler for each command
	commands := []tb.Command{
		{Text: "/start", Description: fmt.Sprintf("Activate the receptions of %s tweets.", ProjectConfig.Hashtag)},
		{Text: "/stop", Description: fmt.Sprintf("Desactivate the receptions of %s tweets.", ProjectConfig.Hashtag)},
		{Text: "/help", Description: "Display all bot commands"},
	}
	myBot.SetCommands(commands)

	myBot.Handle("/start", func(m *tb.Message) {
		var message string

		//We try to fetch if that chat_id is already in the suscriber list and if the status of receiving tweet is ok or not
		//If not we add it with the status ok
		//If yes we say him(her) he already suscribe to receive tweet

		suscriberAlreadyExist := false
		changeOccur := false
		for index, _ := range ProjectConfig.Suscribers {
			if ProjectConfig.Suscribers[index].ChatId == m.Chat.ID {
				suscriberAlreadyExist = true

				if ProjectConfig.Suscribers[index].Status == false {
					ProjectConfig.Suscribers[index].Status = true
					changeOccur = true
					message = fmt.Sprintf("Hey %s ! \nYour status have been change you are now set to receive all %s tweets. \nBy @iamdipanda", m.Chat.Username, ProjectConfig.Hashtag)
					break
				}

				message = fmt.Sprintf("Hey %s ! \nYou are already set to receive all %s tweets. \nBy @iamdipanda", m.Chat.Username, ProjectConfig.Hashtag)
				break
			}
		}

		if suscriberAlreadyExist == false {
			message = fmt.Sprintf("Welcome! \nYou can now receive all %s tweets. \nBy @iamdipanda", ProjectConfig.Hashtag)
			mySuscriber := suscriber.Suscriber{
				UserName: m.Chat.Username,
				ChatId:   m.Chat.ID,
				Status:   true,
			}

			// We add the new user and
			ProjectConfig.Suscribers = append(ProjectConfig.Suscribers, mySuscriber)
			changeOccur = true

			log.Println("#### New Suscriber ####")
			log.Printf("%s is the %v suscriber to join.\n", m.Chat.FirstName, len(ProjectConfig.Suscribers))
		}

		if changeOccur == true {
			// Save it into the config file
			jsonData, err := json.Marshal(ProjectConfig)
			if err == nil {
				ioutil.WriteFile("config.json", jsonData, os.ModePerm)
			}
		}

		// We send the message
		myBot.Send(m.Sender, message)
	})

	myBot.Handle("/stop", func(m *tb.Message) {

		var message string
		suscriberAlreadyExist := false
		changeOccur := false

		for index, _ := range ProjectConfig.Suscribers {
			if ProjectConfig.Suscribers[index].ChatId == m.Chat.ID {
				suscriberAlreadyExist = true

				if ProjectConfig.Suscribers[index].Status == true {
					ProjectConfig.Suscribers[index].Status = false
					changeOccur = true
					message = fmt.Sprintf("Hey %s ! \nYour status has change you will not receive all %s tweets. \nBy @iamdipanda", m.Chat.Username, ProjectConfig.Hashtag)
					break
				}

				message = fmt.Sprintf("Hey %s ! \nYour status has already set to not receive all %s tweets. \nBy @iamdipanda", m.Chat.Username, ProjectConfig.Hashtag)
				break
			}
		}

		if suscriberAlreadyExist == false {
			message = fmt.Sprintf("Hey %s ! \n You have not yet enable the reception of all %s tweets. Send the /start command to enable it.  \n By @iamdipanda", m.Chat.Username, ProjectConfig.Hashtag)
		}

		if changeOccur == true {
			// Save it into the config file
			jsonData, err := json.Marshal(ProjectConfig)
			if err == nil {
				ioutil.WriteFile("config.json", jsonData, 0777)
			}
		}

		// We send the message
		myBot.Send(m.Sender, message)
	})

	myBot.Handle("/help", func(m *tb.Message) {
		message := "Welcome to Bot Help \n \n"

		for _, command := range commands {
			message += fmt.Sprintf("- %s : %s \n", command.Text, command.Description)
		}

		message += "By @iamdipanda"

		// We send the message
		myBot.Send(m.Sender, message)
	})

	log.Println("#### Bot start sucessfully ####")
	myBot.Start()
}

func SendPhoto(suscriber suscriber.Suscriber, photo string, caption string) bool {
	// Will upload the file from disk and send it to suscriber
	fileDir, _ := os.Getwd()
	filePath := path.Join(fileDir, photo)

	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("photo", filepath.Base(file.Name()))
	if err != nil {
		return false
	}
	io.Copy(part, file)
	writer.WriteField("chat_id", strconv.FormatInt(suscriber.ChatId, 10))
	writer.WriteField("caption", caption)
	writer.Close()

	request, _ := http.NewRequest(http.MethodPost,
		"https://api.telegram.org/bot"+ProjectConfig.BotToken+"/sendPhoto",
		body)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return false
	}
	b := response.StatusCode == 200
	return b
}

func SendMessage(subscriber suscriber.Suscriber, message string) bool {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("chat_id", strconv.FormatInt(subscriber.ChatId, 10))
	writer.WriteField("text", message)
	writer.Close()

	request, _ := http.NewRequest(http.MethodPost,
		"https://api.telegram.org/bot"+ProjectConfig.BotToken+"/sendMessage",
		body)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return false
	}
	b := response.StatusCode == 200
	return b
}

func BulkSendPhoto(TweetScreenshot string, tweetLink string, tweetUserName string) {

	caption := fmt.Sprintf("New Tweet from : %s .\n\n-------\n\nLink: %s ", tweetUserName, tweetLink)
	for _, s := range ProjectConfig.Suscribers {
		if s.Status == true {
			SendPhoto(s, TweetScreenshot, caption)
		}
	}

}

func BulkSendText(description string, tweetLink string, tweetUserName string) {

	msg := fmt.Sprintf("New Tweet from : %s .\n\n-------\n\n %s \n\n-------\n\nLink: %s ", tweetUserName, description, tweetLink)
	for _, s := range ProjectConfig.Suscribers {
		if s.Status == true {
			SendMessage(s, msg)
		}
	}

}
