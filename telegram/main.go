package telegram

import (
	"encoding/json"
	//	"fmt"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/naporta/naporta-api/db"
)

var state State

func Start(token string, admins []int, mongo db.Connection) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Printf("Error on get updates from telegram.")
		return
	}

	for update := range updates {
		if update.Message == nil {
			continue // ignore any non-Message Updates
		}

		var isAdmin bool
		for _, item := range admins {
			if item == update.Message.From.ID {
				isAdmin = true
			}
		}

		if !isAdmin {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Opa, parece que você não é admin!")
			bot.Send(msg)
			continue
		}

		if update.Message.IsCommand() && state.ID == HOME {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "verificar":
				arg := update.Message.CommandArguments()

				vendedor, err := mongo.FindOneFalse(arg)
				if err != nil {
					log.Printf("Error: %s", err)
					msg.Text = err.Error()
					bot.Send(msg)
					continue
				}

				state.SetData(update.Message.From.ID, vendedor.ID)
				state.ID = VERIFICAR

				pretty, _ := json.MarshalIndent(vendedor, "", "  ")
				msg.Text = string(pretty)
				bot.Send(msg)
				msg.Text = "digite `rm` para remover, `yes` para autorizar e `no` para deixar como esta."
				bot.Send(msg)
			case "listar":
				vendedores, err := mongo.FindAll("", "")
				if err != nil {
					log.Printf("Error: %s", err)
					msg.Text = err.Error()
					bot.Send(msg)
					continue
				}
				pretty, _ := json.MarshalIndent(vendedores, "", "  ")
				msg.Text = string(pretty)
				bot.Send(msg)

			case "help":
				msg := tgbotapi.NewMessage(
					update.Message.Chat.ID, "Help message!",
				)
				bot.Send(msg)
			}
		}

		if state.ID == VERIFICAR {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			mongoID, exists := state.GetData(update.Message.From.ID)
			if !exists {
				msg.Text = "Deu ruim, comece denovo!"
				bot.Send(msg)
				state.ClearData()
				state.ID = HOME
			}

			switch update.Message.Text {
			case "rm":
				result, err := mongo.Delete(mongoID)
				if err != nil {
					log.Printf("Error: %s", err)
					msg.Text = err.Error()
					bot.Send(msg)
					state.ClearData()
					state.ID = HOME
					continue
				}
				pretty, _ := json.MarshalIndent(result, "", "  ")
				msg.Text = string(pretty)
				bot.Send(msg)
				state.ClearData()
				state.ID = HOME

			case "yes":
				result, err := mongo.UpdateVerificado(mongoID)
				if err != nil {
					log.Printf("Error: %s", err)
					msg.Text = err.Error()
					bot.Send(msg)
					state.ClearData()
					state.ID = HOME
					continue
				}
				log.Printf("Deu certo: ", result)
				msg.Text = "Atualizado com sucesso!"
				bot.Send(msg)
				state.ClearData()
				state.ID = HOME

			case "no":
				state.ClearData()
				state.ID = HOME
				msg.Text = "Vendedor não atualizado."
				bot.Send(msg)
				state.ClearData()
				state.ID = HOME

			}
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	}
}
