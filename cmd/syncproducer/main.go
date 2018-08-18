package main

import (
	"flag"
	"fmt"

	"github.com/Shopify/sarama"
)

var (
	topic = flag.String("topic", "mytopic", "kafka topic you want to publish message to")
	key   = flag.String("key", "mykey", "message key")
	msg   = flag.String("msg", "my message", "the message content you want to publish")
)

func main() {
	flag.Parse()
	cfg := sarama.NewConfig()
	cfg.ClientID = "my-kafka-producer"
	// The following three settings are quite important for Sync Producer
	// it tell the library need to tell us whether the message publish is success or failure
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	// WaitForAll means we need to wait until in sync replica to ack before it return, has strong consistency,but slow
	cfg.Producer.RequiredAcks = sarama.WaitForAll

	// Since we are running a one node
	// Make sure you pass in multiple brokers when it is multiple node cluster
	// client use the address to boostrap , will send request to the broker for meta data
	// and it will discover other brokers from there.
	c, err := sarama.NewClient([]string{"127.0.0.1:9092"}, cfg)
	if nil != err {
		panic(err)
	}

	p, err := sarama.NewSyncProducerFromClient(c)
	if nil != err {
		panic(err)
	}
	defer func() {
		// Please do remember to close the producer, otherwise you are going to have a bad day
		if err := p.Close(); nil != err {
			fmt.Printf("error while closing producer:%s", err)
		}
		if err := c.Close(); nil != err {
			fmt.Printf("error while closing client:%s", err)
		}
	}()

	pmsg := &sarama.ProducerMessage{
		Topic: *topic,
		// Available Encoder include StringEncoder, ByteEncoder
		// Key will be used to decide which partition to publish to,(guarantee order), choose it wisely
		Key: sarama.StringEncoder(*key),
		// if the value is string , then use StringEncoder, otherwise use ByteEncoder
		Value: sarama.StringEncoder(*msg),
	}
	partition, offset, err := p.SendMessage(pmsg)
	if nil != err {
		panic(err)
	}
	fmt.Printf("msg published to partition:%d,offset:%d", partition, offset)
}
