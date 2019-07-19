package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/aravindvs/sarama-cluster"

	"github.com/herlegs/KafkaPlay/constants"
)

const (
	topicName = "xiao"
	groupID   = "tGroup"
)

func main() {
	sarama.Logger = log.New(os.Stdout, "[Sarama] ", log.LstdFlags)

	go Worker("A")
	go Worker("B")
	go Worker("C")

	<-time.After(time.Minute * 10)
}

func Worker(id string) {
	consumerCfg := cluster.NewConfig()
	consumerCfg.ClientID = id
	consumerCfg.Config.Consumer.Offsets.Initial = -2
	consumerCfg.Metadata.RefreshFrequency = time.Second * 10
	consumer, err := cluster.NewConsumer(constants.KafkaCluster, groupID, []string{topicName}, consumerCfg)
	if err != nil {
		fmt.Printf("error creating consumer: %v\n", err)
		return
	}

	msgChan := consumer.Messages()
	if err != nil {
		fmt.Printf("error consumer from partition: %v\n", err)
		return
	}
	for {
		msg, ok := <-msgChan
		if !ok {
			fmt.Printf("channel closed\n")
			return
		}
		fmt.Printf("[%v]-[%v] from Partition[%v]/[%v]\n\n", id, string(msg.Value), msg.Partition, msg.Offset)
		consumer.MarkOffset(msg, "")
	}
}
