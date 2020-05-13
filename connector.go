package main

import (
	"fmt"
	"os"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func ConnectKafka(topicsConfig map[string]Topics, fp func(*kafka.Message, Topics)) {
	broker := os.Getenv("KAFKA_BROKER")
	group := os.Getenv("KAFKA_GROUP")

	if broker == "" {
		broker = "localhost:9092"
	}

	if group == "" {
		group = "test"
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     broker,
		"broker.address.family": "v4",
		"group.id":              group,
		"session.timeout.ms":    6000,
		"auto.offset.reset":     "latest"})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	var topics []string
	for _, element := range topicsConfig {
		topics = append(topics, element.Topic)
	}

	err = c.SubscribeTopics(topics, nil)

	defer c.Close()

	for {
		ev := c.Poll(100)
		if ev == nil {
			continue
		}

		switch e := ev.(type) {
		case *kafka.Message:
			fmt.Println("Now serving!", len(topicsConfig[*e.TopicPartition.Topic].Enhancers))
			fp(e, topicsConfig[*e.TopicPartition.Topic])
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
		default:
			fmt.Printf("Ignored %v\n", e)
		}
	}
}
