package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/Shopify/sarama"
)

func main() {
	topic := flag.String("topic", "mytopic", "kafka topic you want to publish message to")
	flag.Parse()
	cfg := sarama.NewConfig()
	cfg.ClientID = "my-kafka-producer"
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Consumer.Return.Errors = true
	//cfg.Version = sarama.V1_1_0_0
	c, err := sarama.NewClient([]string{"localhost:9092"}, cfg)
	if nil != err {
		panic(err)
	}
	consumer, err := sarama.NewConsumerFromClient(c)
	if nil != err {
		panic(err)
	}
	partitions, err := consumer.Partitions(*topic)
	if nil != err {
		panic(err)
	}
	wg := &sync.WaitGroup{}
	donechan := make(chan struct{})
	for _, p := range partitions {
		pc, err := consumer.ConsumePartition(*topic, p, 0)
		if nil != err {
			panic(err)
		}
		wg.Add(1)
		go processMessage(wg, pc, donechan)
	}
	// clean exit
	close(donechan)
	wg.Wait()
}

func processMessage(wg *sync.WaitGroup, pc sarama.PartitionConsumer, donechan chan struct{}) {
	defer wg.Done()
	for {
		select {
		case err := <-pc.Errors():
			fmt.Printf("Opps, there is an err:%s", err)
		case msg, more := <-pc.Messages():
			if !more {
				return
			}
			fmt.Printf("we got a message,key:%s,msg:%s", string(msg.Key), string(msg.Value))
		case <-donechan:
			pc.AsyncClose()
		}
	}
}
