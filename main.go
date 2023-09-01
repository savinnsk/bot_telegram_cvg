package main

import (
	"encoding/csv"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"os"
	"time"
)

type Contract struct {
	Number         string
	Description    string
	StartDate      time.Time
	ExpirationDate time.Time
}

func parseContractsFromCSV(filePath string) ([]Contract, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skip the header line
	_, _ = reader.Read()

	var contracts []Contract

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		startDate, _ := time.Parse("2/1/2006", record[2])
		expirationDate, _ := time.Parse("2/1/2006", record[3])

		contract := Contract{
			Number:         record[0],
			Description:    record[1],
			StartDate:      startDate,
			ExpirationDate: expirationDate,
		}

		contracts = append(contracts, contract)
	}

	return contracts, nil
}

func handleCommand(command string, contracts []Contract, chatID int64, bot *tgbotapi.BotAPI) error {
	if command == "/check_expired" {
		var expiredContracts []Contract
		now := time.Now()

		for _, contract := range contracts {
			if contract.ExpirationDate.Before(now) {
				expiredContracts = append(expiredContracts, contract)
			}
		}

		if len(expiredContracts) > 0 {
			message := "Expired contracts:\n"
			for _, contract := range expiredContracts {
				message += fmt.Sprintf(
					"Contract: %s\nDescription: %s\nExpiration Date: %s\n\n",
					contract.Number,
					contract.Description,
					contract.ExpirationDate.Format("2006-01-02"),
				)
			}
			msg := tgbotapi.NewMessage(chatID, message)
			_, err := bot.Send(msg)
			return err
		}

		msg := tgbotapi.NewMessage(chatID, "No contracts have expired.")
		_, err := bot.Send(msg)
		return err
	} else if command == "/start" {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Você foi conectado seu id é : %d", chatID))
		_, err := bot.Send(msg)
		return err
	}

	msg := tgbotapi.NewMessage(chatID, "Unknown command.")
	_, err := bot.Send(msg)
	return err
}

func checkExpiringContracts(contracts []Contract) []Contract {
	now := time.Now()
	var expiringContracts []Contract

	for _, contract := range contracts {
		if contract.ExpirationDate.Sub(now).Hours() <= 24*90 {
			expiringContracts = append(expiringContracts, contract)
		}
	}

	return expiringContracts
}

func main() {
	bot, err := tgbotapi.NewBotAPI("6609285497:AAHEUgqCRjrn6_3WzplFFgGJ64VgDPuzsDk")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	contracts, err := parseContractsFromCSV("./database.csv")
	if err != nil {
		log.Fatal(err)
	}

	// Set the time to send messages (e.g., 9:00 AM)
	targetHour := 0
	targetMinute := 1

	for {
		// Get the current time
		now := time.Now()

		// Calculate the time of the next target hour
		nextTargetTime := now.Add(time.Duration(targetHour)*time.Hour + time.Duration(targetMinute)*time.Minute)
		if now.After(nextTargetTime) {
			nextTargetTime = nextTargetTime.Add(24 * time.Hour)
		}

		// Calculate the duration until the next target time
		initialDelay := nextTargetTime.Sub(now)

		// Wait for the initial delay
		time.Sleep(initialDelay)

		expiringContracts := checkExpiringContracts(contracts)
		if len(expiringContracts) > 0 {
			message := "Contracts expiring in 90 days:\n"
			for _, contract := range expiringContracts {
				message += fmt.Sprintf(
					"Contract: %s\nDescription: %s\nExpiration Date: %s\n\n",
					contract.Number,
					contract.Description,
					contract.ExpirationDate.Format("2006-01-02"),
				)
			}

			// Send the warning message to your desired chat ID
			// Replace CHAT_ID with the actual chat ID
			chatID := int64(6252881817)
			warningMsg := tgbotapi.NewMessage(chatID, message)
			_, err := bot.Send(warningMsg)
			if err != nil {
				log.Println("Error sending warning:", err)
			}
		}

	}

}
