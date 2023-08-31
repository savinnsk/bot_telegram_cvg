package entity

import (
	"time"
)

type Contract struct {
	Number         string
	Description    string
	StartDate      time.Time
	ExpirationDate time.Time
}

type ContractSlice []Contract

func (contracts ContractSlice) CheckExpiringContracts() (ContractSlice, error) {
	now := time.Now()
	var expiringContracts ContractSlice

	for _, contract := range contracts {
		if contract.ExpirationDate.Sub(now).Hours() <= 24*90 {
			expiringContracts = append(expiringContracts, contract)
		}
	}

	return expiringContracts, nil
}
