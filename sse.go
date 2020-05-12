package main

import (
	"fmt"
	"log"
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

func formatSSE(event string, data string) []byte {
	eventPayload := "event: " + event + "\n"
    dataLines := strings.Split(data, "\n")
    for _, line := range dataLines {
        eventPayload = eventPayload + "data: " + line + "\n"
    }
    return []byte(eventPayload + "\n")
}

func (broker *Broker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	flusher, ok := rw.(http.Flusher)
	accountNumber := req.URL.Query().Get("account_number")
	room := req.URL.Query().Get("room")

	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	messageChan := make(chan SSEMessage)

	broker.newClients <- messageChan

	defer func() {
		broker.closingClients <- messageChan
	}()

	notify := rw.(http.CloseNotifier).CloseNotify()

	go func() {
		<-notify
		broker.closingClients <- messageChan
	}()

	for {
		// Write to the ResponseWriter
		// Server Sent Events compatible
		channel := <- messageChan
		if accountNumber == channel.accountNumber && room == channel.room {
			fmt.Fprintf(rw, "%s\n", channel.msg)
		}
		flusher.Flush()
	}
}

// Broker factory
func NewServer() (broker *Broker) {
	broker = &Broker{
	  Notifier:       make(chan SSEMessage, 1),
	  newClients:     make(chan chan SSEMessage),
	  closingClients: make(chan chan SSEMessage),
	  clients:        make(map[chan SSEMessage]bool),
	}
  
	go broker.listen()
  
	return
}

func (broker *Broker) listen() {
	for {
	  select {
		case s := <-broker.newClients:
		  // A new client has connected.
		  // Register their message channel
		  broker.clients[s] = true
		  log.Printf("Client added. %d registered clients", len(broker.clients))
  
		case s := <-broker.closingClients:
		  // A client has dettached and we want to
		  // stop sending them messages.
		  delete(broker.clients, s)
		  log.Printf("Removed client. %d registered clients", len(broker.clients))
  
		case event := <-broker.Notifier:
		  // We got a new event from the outside!
		  // Send event to all connected clients
		  for clientMessageChan, _ := range broker.clients {
			clientMessageChan <- event
		  }
	  }
	}
}
