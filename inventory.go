package main

import (
	"encoding/json"
	"fmt"
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

	fmt.Println("Got new inventory event!", msg)
	fmt.Println("Using account", accountNumber)

	if formatted.InventoryHost.Account == accountNumber {
		return true
	}

	return false
}
