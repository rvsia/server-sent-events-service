package main

import (
	"fmt"
)

// SourcesEnhancer is used to determine if the event can be emmited
func SourcesEnhancer(msg string, accountNumber string) bool {

	fmt.Println("Got new sources event", msg)
	fmt.Println("Using account", accountNumber)

	return true
}
