package main

import (
	"encoding/json"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/my1562/telegram-broadcaster/config"
	"github.com/my1562/telegram-broadcaster/tg"
	"github.com/streadway/amqp"
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

func main() {

	c := dig.New()
	c.Provide(config.NewConfig)
	c.Provide(tg.NewTg)

	err := c.Invoke(func(config *config.Config, bot *tgbotapi.BotAPI) {

		conn, err := amqp.Dial(config.RabbitmqURL)
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()

		q, err := ch.QueueDeclare(
			"hello", // name
			false,   // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		failOnError(err, "Failed to declare a queue")

		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		failOnError(err, "Failed to register a consumer")
		forever := make(chan bool)

		go func() {
			for d := range msgs {
				log.Printf("Received a message: %s", d.Body)

				notification := &TelegramNotification{}
				json.Unmarshal(d.Body, notification)
				if notification.ChatID == 0 {
					log.Println("[warning] ChatID==0")
					continue
				}
				if notification.Message == "" {
					log.Println("[warning] Message is empty")
					continue
				}

				msg := tgbotapi.NewMessage(notification.ChatID, notification.Message)
				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			}
		}()

		log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

		<-forever

	})
	if err != nil {
		panic(err)
	}
}
