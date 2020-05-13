package main

import (
	"encoding/json"
)

type Host struct {
	Account string `json:"account"`
}

type InventoryMessage struct {
	InventoryHost Host `json:"host"`
}

func InventoryEnhancer(msg string, accountNumber string) bool {
	var formatted InventoryMessage
	json.Unmarshal([]byte(msg), &formatted)

	if formatted.InventoryHost.Account == accountNumber {
		return true
	}

	return false
}
