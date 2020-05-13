package enhancers

import (
	"encoding/json"
	"fmt"
)

type Host struct {
	Account string `json:"account"`
}

type InventoryMessage struct {
	Host Host `json:"host"`
}

func InventoryEnhancer(msg string, accountNumber string) bool {
	var formatted InventoryMessage
	json.Unmarshal([]byte(msg), &formatted)
	fmt.Println("inventory data", msg)

	if formatted.Host.Account == accountNumber {
		return true
	}

	return false
}
