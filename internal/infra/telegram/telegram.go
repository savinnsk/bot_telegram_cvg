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

func InitConnection(apiToken string) (*Telegram, error) {

	connection, err := tgbotapi.NewBotAPI(apiToken)

	if err != nil {
		fmt.Println("error at init connection")
	}
	return &Telegram{
		Connection: connection,
	}, err

}

func (t *Telegram) SendMessage(message string, chatID int64) error {
	newMessage := tgbotapi.NewMessage(chatID, message)
	_, err := t.Connection.Send(newMessage)

	return err

}

func (t *Telegram) HandleCommand(command string, contracts []entity.Contract, chatID int64) error {
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
			err := t.SendMessage(message, chatID)

			return err
		}

		err := t.SendMessage("Não Há contratos expirados.", chatID)

		return err
	} else if command == "/start" {
		err := t.SendMessage(fmt.Sprintf("Você foi conectado seu id é : %d", chatID), chatID)
		return err
	}

	err := t.SendMessage("Comando não encontrado.", chatID)

	return err
}
