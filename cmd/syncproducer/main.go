package main

import (
	"flag"
	"fmt"

	"github.com/Shopify/sarama"
)

func main() {
	topic := flag.String("topic", "mytopic", "kafka topic you want to publish message to")
	key := flag.String("key", "mykey", "message key")
	msg := flag.String("msg", "my message", "the message content you want to publish")
	flag.Parse()
	cfg := sarama.NewConfig()
	cfg.ClientID = "my-kafka-producer"
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	//cfg.Version = sarama.V1_1_0_0
	c, err := sarama.NewClient([]string{"localhost:9092"}, cfg)
	if nil != err {
		panic(err)
	}
	p, err := sarama.NewSyncProducerFromClient(c)
	if nil != err {
		panic(err)
	}
	defer func() {
		if err := p.Close(); nil != err {
			fmt.Printf("error while closing producer:%s", err)
		}
	}()
	pmsg := &sarama.ProducerMessage{
		Topic: *topic,
		Key:   sarama.StringEncoder(*key),
		Value: sarama.StringEncoder(*msg),
	}
	partition, offset, err := p.SendMessage(pmsg)
	if nil != err {
		panic(err)
	}
	fmt.Printf("msg published to partition:%d,offset:%d", partition, offset)
}
