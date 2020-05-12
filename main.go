package main

import (
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"github.com/joho/godotenv"
	"encoding/json"
	"io/ioutil"
)

type Topics struct {
	Topic string `json:"topic"`
}

func readTopics() []Topics {
	file, _ := ioutil.ReadFile("topics.json")
	data := make([]Topics,0)

	_ = json.Unmarshal([]byte(file), &data)

	return data
}

func printKafka(msg *kafka.Message) {
	fmt.Printf("%% Message on %s:\n%s\n", msg.TopicPartition, string(msg.Value))
}

func main() {
	topicsConfig := readTopics()
	_ = godotenv.Load()
	
	var topics []string
	for i := 0; i < len(topicsConfig); i++ {
		topics = append(topics, topicsConfig[i].Topic)
	}

	connectKafka(topics, printKafka)
}
