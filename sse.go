package main

import (
	"net/http"
	"strings"
)

type SSEMessage struct {
	room string
	accountNumber string
	msg []byte
}

type Broker struct {
    Notifier chan SSEMessage
    newClients chan chan SSEMessage
    closingClients chan chan SSEMessage
    clients map[chan SSEMessage]bool
}

var messageChannels = make(map[chan SSEMessage]bool)

func formatSSE(event string, data string) []byte {
	eventPayload := "event: " + event + "\n"
    dataLines := strings.Split(data, "\n")
    for _, line := range dataLines {
        eventPayload = eventPayload + "data: " + line + "\n"
    }
    return []byte(eventPayload + "\n")
}

func listenHandler(w http.ResponseWriter, req *http.Request) {
	accountNumber := req.URL.Query().Get("account_number")
	room := req.URL.Query().Get("room")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Access-Control-Allow-Origin", "*")

    _messageChannel := make(chan SSEMessage)
    messageChannels[_messageChannel] = true

    for {
        select {
		case channel := <-_messageChannel:
			if accountNumber == channel.accountNumber && room == channel.room {
				w.Write(channel.msg)
			}
            w.(http.Flusher).Flush()
        case <-req.Context().Done():
            delete(messageChannels, _messageChannel)
            return
        }
    }
}
