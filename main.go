package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type Topics struct {
	Topic string `json:"topic"`
	Room  string `json:"room"`
}

func readTopics() map[string]Topics {
	file, _ := ioutil.ReadFile("topics.json")
	data := make([]Topics, 0)

	json.Unmarshal([]byte(file), &data)

	dataMap := make(map[string]Topics)
	for i := 0; i < len(data); i++ {
		dataMap[data[i].Topic] = data[i]
	}

	return dataMap
}

func sendToListener(msg *kafka.Message, topic string) {
	accountNumber := "55"

	go func() {
		for messageChannel, connectorInfo := range messageChannels {
			msg := formatSSE("testing", string(msg.Value))
			if connectorInfo.accountNumber == accountNumber {
				fmt.Println("User found")
				if topic == "" {
					fmt.Println("No room, broadcasting!")
					messageChannel <- msg
				} else if connectorInfo.room == topic {
					fmt.Println("Sending to specific room")
					messageChannel <- msg
				} else {
					fmt.Println("Not sending", connectorInfo.room, topic)
				}
			}
		}
	}()
	fmt.Printf("%% Message on %s:\n%s\n", msg.TopicPartition, string(msg.Value))
}

func main() {
	topicsConfig := readTopics()
	godotenv.Load()
	apiVersion := os.Getenv("API_VERSION")

	if apiVersion == "" {
		apiVersion = "v1"
	}

	go func() {
		connectKafka(topicsConfig, sendToListener)
	}()

	// go func() {
	// 	for {
	// 		time.Sleep(time.Second * 2)
	// 		log.Println("Receiving event")
	// 		sendToListener(fmt.Sprintf("the time is %v", time.Now()), "test")
	// 	}
	// }()

	http.HandleFunc(fmt.Sprintf("/api/notifier/%s/connect", apiVersion), listenHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
