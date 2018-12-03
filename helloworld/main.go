package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"github.com/herlegs/KafkaPlay/constants"
	"time"
)

func main() {
	producerCfg := sarama.NewConfig()
	producerCfg.Producer.Return.Successes = true
	producer,err := sarama.NewSyncProducer(constants.KafkaCluster, producerCfg)
	if err != nil {
		fmt.Printf("error creating producer: %v\n",err)
		return
	}

	// read consumer offsets
	go func(){
		cgCfg := cluster.NewConfig()
		cgCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

		clusterConsumer,err := cluster.NewConsumer(constants.KafkaCluster, "testGroup", []string{"__consumer_offsets"}, cgCfg)
		if err != nil {
			fmt.Printf("error create offset consumer: %v\n",err)
			return
		} else {
			fmt.Printf("created offset consumer!\n",)
		}

		dataCh := clusterConsumer.Messages()

		i := 0
		for msg := range dataCh {
			fmt.Printf("[offsetconsumer]topic[%v]key[%v]value[%v]offset[%v]\n",msg.Topic,string(msg.Key),string(msg.Value),msg.Offset)
			i++
		}
	}()

	//sarama consumer
	//go func(){
	//	pConsumer, err := consumer.ConsumePartition("testTopic", 0, sarama.OffsetOldest)
	//	if err != nil {
	//		fmt.Printf("error create sarama consumer: %v\n",err)
	//		return
	//	}
	//
	//	dataCh := pConsumer.Messages()
	//	for msg := range dataCh {
	//		fmt.Printf("topic[%v]key[%v]value[%v]offset[%v]\n",msg.Topic,string(msg.Key),string(msg.Value),msg.Offset)
	//	}
	//}()


	//sarama-cluster consumer
	go func(){
		cgCfg := cluster.NewConfig()
		cgCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

		clusterConsumer,err := cluster.NewConsumer(constants.KafkaCluster, "testGroup", []string{"testTopic"}, cgCfg)
		if err != nil {
			fmt.Printf("error create sarama-cluster consumer: %v\n",err)
			return
		} else {
			fmt.Printf("created sarama-cluster consumer!\n",)
		}

		dataCh := clusterConsumer.Messages()

		i := 0
		for msg := range dataCh {
			fmt.Printf("[sarama-cluster consumer]topic[%v]key[%v]value[%v]offset[%v]\n",msg.Topic,string(msg.Key),string(msg.Value),msg.Offset)
			i++
		}
	}()


	for {
		producer.SendMessage(&sarama.ProducerMessage{
			Topic:"testTopic",
			Key: sarama.StringEncoder("testKey"),
			Value: sarama.StringEncoder("testVal"),
		})
		_=producer
		time.Sleep(time.Second*10)
	}
}
