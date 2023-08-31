package infra

import (
	"fmt"
	"time"

	"github.com.br/telegram_go_bot/internal/entity/contract"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram struct {
	Connection *tgbotapi.BotAPI
}

func NewConnection(connection *tgbotapi.BotAPI) *Telegram {
	return &Telegram{
		Connection: connection,
	}
}

func SendMessage(bot *Telegram, message string, chatID int64) {
	newMessage := tgbotapi.NewMessage(chatID, message)
	_, err := bot.Connection.Send(newMessage)

	if err != nil {
		fmt.Println("error at send message")
	}
}

func HandleCommand(command string, contracts []entity.Contract, chatID int64, bot *tgbotapi.BotAPI) error {
	if command == "/check_expired" {
		var expiredContracts []entity.Contract
		now := time.Now()

		for _, contract := range contracts {
			if contract.ExpirationDate.Before(now) {
				expiredContracts = append(expiredContracts, contract)
			}
		}

		if len(expiredContracts) > 0 {
			message := "Contratos expirados:\n"
			for _, contract := range expiredContracts {
				message += fmt.Sprintf(
					"Contrato: %s\nDescrição: %s\nData de Expiração: %s\n\n",
					contract.Number,
					contract.Description,
					contract.ExpirationDate.Format("2006-01-02"),
				)
			}
			msg := tgbotapi.NewMessage(chatID, message)
			_, err := bot.Send(msg)
			return err
		}

		msg := tgbotapi.NewMessage(chatID, "Não Há contratos expirados.")
		_, err := bot.Send(msg)
		return err
	} else if command == "/start" {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Você foi conectado seu id é : %d", chatID))
		_, err := bot.Send(msg)
		return err
	}

	msg := tgbotapi.NewMessage(chatID, "Comando não encontrado.")
	_, err := bot.Send(msg)
	return err
}
