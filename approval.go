package main

import (
	"encoding/json"
	"fmt"
)

func ApprovalEnhancer(msg string, accountNumber string) bool {

	fmt.Println("Got new approval event!", msg)
	fmt.Println("Using account", accountNumber)

	return false
}
