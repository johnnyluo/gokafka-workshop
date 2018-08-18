package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
)

func main() {
	topic := flag.String("topic", "mytopic", "kafka topic you want to publish message to")
	flag.Parse()
	cfg := sarama.NewConfig()
	cfg.ClientID = "my-kafka-consumer"
	cfg.Consumer.Return.Errors = true

	c, err := sarama.NewClient([]string{"localhost:9092"}, cfg)
	if nil != err {
		panic(err)
	}
	consumer, err := sarama.NewConsumerFromClient(c)
	if nil != err {
		panic(err)
	}
	defer func() {
		if err := consumer.Close(); nil != err {
			fmt.Printf("fail to close consumer,err:%s\n", err)
		}
		if err := c.Close(); nil != err {
			fmt.Printf("fail to close the client,err:%s\n", err)
		}
	}()
	partitions, err := consumer.Partitions(*topic)
	if nil != err {
		panic(err)
	}
	wg := &sync.WaitGroup{}
	done := make(chan struct{})
	for _, p := range partitions {
		// start from offset 0
		pc, err := consumer.ConsumePartition(*topic, p, 0)
		if nil != err {
			panic(err)
		}
		wg.Add(1)
		go processMessage(wg, pc, done)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	close(done)
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
			fmt.Printf("we got a message,key:%s,msg:%s partition:%d, offset:%d \n", string(msg.Key), string(msg.Value), msg.Partition, msg.Offset)
		case <-donechan:
			pc.AsyncClose()
		}
	}
}
