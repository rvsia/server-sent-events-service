package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gobuffalo/packr"
	"github.com/joho/godotenv"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type Topics struct {
	Topic     string   `json:"topic"`
	Room      string   `json:"room"`
	Event     string   `json:"event"`
	Enhancers []string `json:"enhancers"`
}

func readTopics() map[string]Topics {
	box := packr.NewBox("./static")
	file, foundError := box.FindString("topics.json")
	if foundError != nil {
		fmt.Println("Error while fetching file")
		file = os.Getenv("CONFIG_JSON")
	}
	data := make([]Topics, 0)

	err := json.Unmarshal([]byte(file), &data)
	if err != nil {
		fmt.Println(err)
	}

	dataMap := make(map[string]Topics)
	for i := 0; i < len(data); i++ {
		dataMap[data[i].Topic] = data[i]
	}

	return dataMap
}

func sendToListener(kafkaMessage *kafka.Message, topic Topics) {
	if topic.Event == "" {
		topic.Event = "notification"
	}

	enhancers := map[string](func(string, string) bool){
		"inventory": InventoryEnhancer,
		"approval":  ApprovalEnhancer,
	}
	go func() {
		for messageChannel, connectorInfo := range MessageChannels {
			canSend := true
			for i := 0; i < len(topic.Enhancers); i++ {
				canSend = enhancers[topic.Enhancers[i]](string(kafkaMessage.Value), connectorInfo.accountNumber)
			}
			if canSend {
				msg := FormatSSE(topic.Event, string(kafkaMessage.Value))
				if topic.Room == "" {
					fmt.Println("No room, broadcasting!")
					messageChannel <- msg
				} else if connectorInfo.room == topic.Room {
					fmt.Println("Sending to specific room")
					messageChannel <- msg
				} else {
					fmt.Println("Not sending", connectorInfo.room, topic.Room)
				}
			}
		}
	}()

	fmt.Printf("%% Message on %s:\n%s\n", kafkaMessage.TopicPartition, string(kafkaMessage.Value))
}

func main() {
	godotenv.Load()
	topicsConfig := readTopics()
	apiVersion := os.Getenv("API_VERSION")
	appName := os.Getenv("APP_NAME")

	if apiVersion == "" {
		apiVersion = "v1"
	}

	if appName == "" {
		appName = "notifier"
	}

	go ConnectKafka(topicsConfig, sendToListener)

	http.HandleFunc(fmt.Sprintf("/api/%s/%s/connect", appName, apiVersion), ListenHandler)
	http.HandleFunc(fmt.Sprintf("/api/%s/%s/lubdub", appName, apiVersion), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "lubdub")
	})

	srv := &http.Server{
		Addr:         ":3000",
		Handler:      nil,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}

	log.Fatal(srv.ListenAndServe())
}
