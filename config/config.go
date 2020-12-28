package config

import (
	"encoding/json"
	"errors"
	"github.com/DipandaAser/tweetwatcher/bot/suscriber"
	"io/ioutil"
)

type Configuration struct {
	Hashtag    string                `json:"Hashtag"`
	BotToken   string                `json:"BotToken"`
	Port       string                `json:"Port"`
	PublicURL  string                `json:"PublicURL"`
	Suscribers []suscriber.Suscriber `json:"Suscribers"`
}

// Check check if requires fields are provided
func (config *Configuration) Check() error {
	if config.Hashtag == "" {
		return errors.New("Please specify a Hashtag in config file.")
	}
	if config.BotToken == "" {
		return errors.New("Please specify a BotToken in config file.")
	}
	if config.Port == "" {
		return errors.New("Please specify a Port in config file.")
	}
	if config.PublicURL == "" {
		return errors.New("Please specify the PublicURL webbhook of the bot in config file.")
	}
	return nil
}

func CreateConfigTemplate() {
	var conf Configuration
	var suscriber suscriber.Suscriber
	conf.Suscribers = append(conf.Suscribers, suscriber)
	jsonData, err := json.Marshal(conf)
	if err == nil {
		ioutil.WriteFile("config-template.json", jsonData, 0777)
	}
}
