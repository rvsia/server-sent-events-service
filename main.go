package main

import (
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"github.com/joho/godotenv"
	"encoding/json"
	"io/ioutil"
	"log"
	// "time"
	"net/http"
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

func sendToListener(broker *Broker) (func(*kafka.Message)) {
	return func (msg *kafka.Message) {
		broker.Notifier <- []byte(msg.Value)
		fmt.Printf("%% Message on %s:\n%s\n", msg.TopicPartition, string(msg.Value))
	}
}

func main() {
	topicsConfig := readTopics()
	_ = godotenv.Load()
	
	var topics []string
	for i := 0; i < len(topicsConfig); i++ {
		topics = append(topics, topicsConfig[i].Topic)
	}

	broker := NewServer()

	// go func() {
	// 	for {
	// 		time.Sleep(time.Second * 2)
	// 		eventString := fmt.Sprintf("the time is %v", time.Now())
	// 		log.Println("Receiving event")
	// 		broker.Notifier <- []byte(eventString)
	// 	}
	// }()

	go func(){
		connectKafka(topics, sendToListener(broker))
	}()

	log.Fatal("HTTP server error: ", http.ListenAndServe("localhost:3000", broker))
}
