package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/my1562/telegram-broadcaster/config"
)

func NewTg(conf *config.Config) (*tgbotapi.BotAPI, error) {
	api, err := tgbotapi.NewBotAPI(conf.TGToken)
	if err != nil {
		return nil, err
	}
	api.Debug = true

	return api, nil
}
