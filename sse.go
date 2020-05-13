package main

import (
	"fmt"
	"net/http"
	"strings"
)

type SSEMessage struct {
	room          string
	accountNumber string
}

var MessageChannels = make(map[chan []byte]SSEMessage)

func FormatSSE(event string, data string) []byte {
	eventPayload := "event: " + event + "\n"
	dataLines := strings.Split(data, "\n")
	for _, line := range dataLines {
		eventPayload = eventPayload + "data: " + line + "\n"
	}
	return []byte(eventPayload + "\n")
}

func ListenHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var sseChannel SSEMessage
	sseChannel.accountNumber = req.URL.Query().Get("account_number")
	sseChannel.room = req.URL.Query().Get("room")

	_messageChannel := make(chan []byte)
	MessageChannels[_messageChannel] = sseChannel

	fmt.Println("We have a new connection!", sseChannel)

	for {
		select {
		case channel := <-_messageChannel:
			w.Write(channel)
			w.(http.Flusher).Flush()
		case <-req.Context().Done():
			delete(MessageChannels, _messageChannel)
			return
		}
	}
}
