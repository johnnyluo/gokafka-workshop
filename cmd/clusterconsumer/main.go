package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

func main() {
	topic := flag.String("topic", "mytopic", "kafka topic you want to publish message to")
	flag.Parse()
	cfg := cluster.NewConfig()
	cfg.ClientID = "my-kafka-cluster-consumer"
	cfg.Group.PartitionStrategy = cluster.StrategyRoundRobin
	cfg.Group.Return.Notifications = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Consumer.Offsets.CommitInterval = 500 * time.Millisecond
	cfg.Consumer.Return.Errors = true
	cfg.Version = sarama.V1_1_0_0
	c, err := cluster.NewClient([]string{"localhost:9092"}, cfg)
	if nil != err {
		panic(err)
	}

	consumer, err := cluster.NewConsumerFromClient(c, "my-kafka-consumer-group", []string{*topic})
	if nil != err {
		panic(err)
	}
	wg := &sync.WaitGroup{}
	donechan := make(chan struct{})
	wg.Add(2)
	go processMessage(wg, consumer, donechan)
	go processNotification(wg, consumer.Notifications())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	close(donechan)
	wg.Wait()
}

func processNotification(wg *sync.WaitGroup, notification <-chan *cluster.Notification) {
	defer wg.Done()
	for n := range notification {
		fmt.Printf("notification type:%s\n", n.Type)
		for topic, partitions := range n.Current {
			fmt.Printf("topic:%s,current partitions:%+v\n", topic, partitions)
		}
		for topic, partitions := range n.Released {
			fmt.Printf("topic:%s,released partitions:%+v\n", topic, partitions)
		}
		for topic, partitions := range n.Claimed {
			fmt.Printf("topic:%s,claimed partitions:%+v\n", topic, partitions)
		}
	}
}
func processMessage(wg *sync.WaitGroup, pc *cluster.Consumer, donechan chan struct{}) {
	defer wg.Done()
	for {
		select {
		case err := <-pc.Errors():
			if nil == err {
				return
			}
			fmt.Printf("Opps, there is an err:%s\n", err)
		case msg, more := <-pc.Messages():
			if !more {
				return
			}
			fmt.Printf("we got a message,key:%s,msg:%s partition:%d, offset:%d \n", string(msg.Key), string(msg.Value), msg.Partition, msg.Offset)
			pc.MarkOffset(msg, "")
		case <-donechan:
			pc.Close()
		}
	}
}
