package entity

import (
	"fmt"
	"time"
)

type Contract struct {
	Number         string
	Description    string
	StartDate      time.Time
	ExpirationDate time.Time
}

func FilterExpiringContracts(contracts []Contract) string {
	var expiringContracts []Contract
	now := time.Now()

	message := "Contracts expiring in 90 days:\n"

	for _, contract := range contracts {
		if contract.ExpirationDate.Sub(now).Hours() <= 24*90 {
			expiringContracts = append(expiringContracts, contract)
			message += fmt.Sprintf(
				"Contract: %s\nDescription: %s\nExpiration Date: %s\n\n",
				contract.Number,
				contract.Description,
				contract.ExpirationDate.Format("2006-01-02"),
			)
		}
	}

	// Check if there are expiring contracts, if not, return a different message
	if len(expiringContracts) == 0 {
		message = "No contracts expiring in 90 days."
	}

	return message
}
