package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hibiken/asynq"
	"github.com/my1562/queue"
	"github.com/my1562/telegram-broadcaster/config"
	"github.com/my1562/telegram-broadcaster/tg"
	"go.uber.org/dig"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type TelegramNotification struct {
	ChatID  int64
	Message string
}

type NotifierImpl struct {
	bot *tgbotapi.BotAPI
}

func NewNotifierImpl(bot *tgbotapi.BotAPI) queue.INotifyExecutor {
	return &NotifierImpl{bot: bot}
}

func (n *NotifierImpl) Notify(chatID int64, message string) error {
	log.Printf("Sending to %d: %s", chatID, message)

	msg := tgbotapi.NewMessage(chatID, message)

	if _, err := n.bot.Send(msg); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func main() {

	c := dig.New()
	c.Provide(config.NewConfig)
	c.Provide(tg.NewTg)
	c.Provide(NewNotifierImpl)
	c.Provide(queue.NewNotifyHandler)

	err := c.Invoke(func(config *config.Config, bot *tgbotapi.BotAPI, handler *queue.NotifyHandler) {

		redis := asynq.RedisClientOpt{Addr: config.Redis}
		server := asynq.NewServer(redis, asynq.Config{
			Concurrency: 1,
		})
		mux := asynq.NewServeMux()
		mux.Handle(queue.TaskTypeNotify, handler)

		if err := server.Run(mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}

	})
	if err != nil {
		panic(err)
	}
}
