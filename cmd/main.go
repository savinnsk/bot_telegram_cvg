package main

import (
	"fmt"
	bot "github.com.br/telegram_go_bot/internal/entity/bot"
	contract "github.com.br/telegram_go_bot/internal/entity/contract"
	cvg "github.com.br/telegram_go_bot/internal/infra/cvg"
	telegram "github.com.br/telegram_go_bot/internal/infra/telegram"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	cron "github.com/robfig/cron/v3"
	"log"
	"strings"
)

func sendMessages(telegramClient *telegram.Telegram, botCredentials *bot.Bot) {
	contracts, err := cvg.ParseContractsFromCSV("../database.csv")

	if err != nil {
		println("contracts not found")
	}

	expiredContracts := contract.FilterExpiringContracts(contracts)

	for _, chatID := range botCredentials.ChatIds {
		err := telegramClient.SendMessage(expiredContracts, chatID)
		if err != nil {
			fmt.Printf("Error sending message to chat %d: %v\n", chatID, err)
		} else {
			fmt.Printf("Message sent successfully to chat %d\n", chatID)
		}
	}
}

func handleIncomingMessages(telegramClient *telegram.Telegram, contracts []contract.Contract) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := telegramClient.Connection.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		command := update.Message.Text
		chatID := update.Message.Chat.ID

		// Handle commands
		if strings.HasPrefix(command, "/") {
			err := telegramClient.HandleCommand(command, contracts, chatID)
			if err != nil {
				fmt.Printf("Error handling command %s: %v\n", command, err)
			}
		}
	}
}

func main() {

	botCredentials, err := bot.NewBot("YOUR_BOT_TOKEN", []int64{1 /*YOUR_CHAT_ID*/})

	if err != nil {
		panic(err)
	}

	telegram, err := telegram.InitConnection(botCredentials.ApiToken)

	if err != nil {
		panic(err)
	}

	telegram.Connection.Debug = true

	log.Printf("Authorized on account %s", telegram.Connection.Self.UserName)

	go func() {
		contracts, err := cvg.ParseContractsFromCSV("../database.csv")

		if err != nil {
			println("contracts not found")
		}

		handleIncomingMessages(telegram, contracts)
	}()

	c := cron.New()

	// Add TIME a cron job to send messages
	_, err = c.AddFunc("* * * * *", func() {
		sendMessages(telegram, botCredentials)
	})

	if err != nil {
		panic(err)
	}
	// Start the cron job scheduler
	c.Start()

	// Keep the program running indefinitely (or until terminated)
	select {}

}
