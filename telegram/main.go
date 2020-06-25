package telegram

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/naporta/naporta-api/db"
)

func Start(token string, mongo db.Connection) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue // ignore any non-Message Updates
		}

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "add":
				msg.Text = fmt.Sprintf("%v", update.Message)
				bot.Send(msg)
				data := strings.Split(update.Message.Text, ",")
				whatsapp, err := strconv.ParseInt(data[4], 10, 64)
				bloco, err := strconv.ParseInt(data[6], 10, 64)
				apt, err := strconv.ParseInt(data[7], 10, 64)
				if err != nil {
					log.Println(err)
				}

				vendedor := db.Vendedor{
					Empresa:     data[1],
					Responsavel: data[2],
					Produtos:    data[3],
					Whatsapp:    whatsapp,
					Condominio:  data[5],
					Bloco:       bloco,
					Apt:         apt,
					Pagamento:   data[8],
				}
				mongo.Insert(vendedor)
				msg.Text = "Deu bom pai!"
				bot.Send(msg)
			case "listar":
				vendedores, err := mongo.FindAll()
				if err != nil {
					log.Printf("Errou: %s", err)
					continue
				}
				msg := tgbotapi.NewMessage(
					update.Message.Chat.ID, fmt.Sprintf("%v", vendedores),
				)
				bot.Send(msg)

			case "help":
				msg := tgbotapi.NewMessage(
					update.Message.Chat.ID, "Help message!",
				)
				bot.Send(msg)
			}
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	}
}
