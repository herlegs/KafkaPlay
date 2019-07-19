package main

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"github.com/herlegs/KafkaPlay/constants"
)

func main() {
	producerCfg := sarama.NewConfig()
	producerCfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(constants.KafkaCluster, producerCfg)
	if err != nil {
		fmt.Printf("error creating producer: %v\n", err)
		return
	}

	consumerCfg := cluster.NewConfig()
	consumer, err := cluster.NewConsumer(constants.KafkaCluster, "oneGroup", []string{"test"}, consumerCfg)
	if err != nil {
		fmt.Printf("error creating consumer: %v\n", err)
		return
	}

	go func() {
		pConsumer := consumer.Messages()
		if err != nil {
			fmt.Printf("error consumer from partition: %v\n", err)
			return
		}
		msgChan := pConsumer.Messages()
		pConsumer.HighWaterMarkOffset()
		for {
			msg, ok := <-msgChan
			if !ok {
				fmt.Printf("channel closed\n")
				return
			}
			fmt.Printf("topic[%v]key[%v]value[%v]offset[%v]\n", msg.Topic, string(msg.Key), string(msg.Value), msg.Offset)
		}
	}()

	for {
		producer.SendMessage(&sarama.ProducerMessage{
			Topic: "testTopic",
			Key:   sarama.StringEncoder("testKey"),
			Value: sarama.StringEncoder("testVal"),
		})

		time.Sleep(time.Second * 10)
	}
}
