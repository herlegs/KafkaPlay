package main

import (
	"fmt"
	"time"

	"github.com/aravindvs/sarama-cluster"

	"github.com/Shopify/sarama"
	"github.com/herlegs/KafkaPlay/constants"
)

const (
	topicName = "test"

	UserName = "alice"
	Password = "alicepw"
)

func main() {
	//sarama.Logger = log.New(os.Stdout, "[Sarama] ", log.LstdFlags)
	producerCfg := sarama.NewConfig()
	//producerCfg.Metadata.RefreshFrequency = time.Second * 20
	producerCfg.Producer.Return.Successes = true
	producerCfg.Net.SASL.Enable = true
	producerCfg.Net.SASL.User = UserName
	producerCfg.Net.SASL.Password = Password

	producer, err := sarama.NewSyncProducer(constants.KafkaClusterSecure, producerCfg)
	if err != nil {
		fmt.Printf("error creating producer: %v\n", err)
		return
	}

	go Worker("A")
	go Worker("B")
	go Worker("C")

	i := 0
	for {
		val := fmt.Sprintf("val%v", i)
		msg := &sarama.ProducerMessage{
			Topic: topicName,
			Key:   sarama.StringEncoder(fmt.Sprintf("key%v", i)),
			Value: sarama.StringEncoder(val),
		}
		partition, _, err := producer.SendMessage(msg)
		if err != nil {
			fmt.Printf("err send msg: %v\n", err)
		}
		fmt.Printf("sent msg %v to partition %v\n", val, partition)
		i++
		time.Sleep(time.Second * 3)
	}
}

func Worker(id string) {
	consumerCfg := cluster.NewConfig()
	consumerCfg.Metadata.RefreshFrequency = time.Second * 10
	consumerCfg.Net.SASL.Enable = true
	consumerCfg.Net.SASL.User = UserName
	consumerCfg.Net.SASL.Password = Password
	consumer, err := cluster.NewConsumer(constants.KafkaClusterSecure, "tGroup", []string{topicName}, consumerCfg)
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
