package main

import (
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"github.com/joho/godotenv"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"net/http"
)

type Topics struct {
	Topic string `json:"topic"`
	Room string `json:"room"`
}

func readTopics() map[string]Topics {
	file, _ := ioutil.ReadFile("topics.json")
	data := make([]Topics,0)

	json.Unmarshal([]byte(file), &data)

	dataMap := make(map[string]Topics)
	for i := 0; i < len(data); i++ {
		dataMap[data[i].Topic] = data[i]
	}

	return dataMap
}

func sendToListener(msg *kafka.Message, topic string) {
	var message SSEMessage
	message.accountNumber = "55"
	message.room = topic
	message.msg = formatSSE("testing", string(msg.Value))

	go func() {
		for messageChannel := range messageChannels {
			messageChannel <- message
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

	go func(){
		connectKafka(topicsConfig, sendToListener)
	}()

	http.HandleFunc(fmt.Sprintf("/api/notifier/%s/connect", apiVersion), listenHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
