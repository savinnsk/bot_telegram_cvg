package infra

import (
	"encoding/csv"
	"github.com.br/telegram_go_bot/internal/entity/contract"
	"io"
	"os"
	"time"
)

func ParseContractsFromCSV(filePath string) ([]entity.Contract, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skip the header line
	_, _ = reader.Read()

	var contracts []entity.Contract

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

		contract := entity.Contract{
			Number:         record[0],
			Description:    record[1],
			StartDate:      startDate,
			ExpirationDate: expirationDate,
		}

		contracts = append(contracts, contract)
	}

	return contracts, nil
}
