package main

import (
	"fmt"
	"log"
	"net/http"
)

type Broker struct {
    Notifier chan []byte
    newClients chan chan []byte
    closingClients chan chan []byte
    clients map[chan []byte]bool
}

func (broker *Broker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	flusher, ok := rw.(http.Flusher)

	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	messageChan := make(chan []byte)

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
		fmt.Fprintf(rw, "data: %s\n\n", <-messageChan)
	  
		// Flush the data immediatly instead of buffering it for later.
		flusher.Flush()
	}
}

// Broker factory
func NewServer() (broker *Broker) {
	broker = &Broker{
	  Notifier:       make(chan []byte, 1),
	  newClients:     make(chan chan []byte),
	  closingClients: make(chan chan []byte),
	  clients:        make(map[chan []byte]bool),
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
