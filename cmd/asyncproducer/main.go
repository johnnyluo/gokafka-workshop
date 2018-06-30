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
)

var (
	topic = flag.String("topic", "mytopic", "kafka topic you want to publish message to")
	num   = flag.Int("msgcount", 10, "how many messages you want to generate")
)

func main() {
	flag.Parse()
	cfg := sarama.NewConfig()
	cfg.ClientID = "my-kafka-producer"
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Flush.Bytes = 65535
	cfg.Producer.Flush.Frequency = 500 * time.Millisecond
	//cfg.Producer.Flush.Messages = 100
	//cfg.Version = sarama.V1_1_0_0
	c, err := sarama.NewClient([]string{"127.0.0.1:9092"}, cfg)
	if nil != err {
		panic(err)
	}
	p, err := sarama.NewAsyncProducerFromClient(c)
	if nil != err {
		panic(err)
	}
	defer func() {
		if err := p.Close(); nil != err {
			fmt.Printf("error while closing producer:%s", err)
		}
	}()
	wg := &sync.WaitGroup{}
	donechan := make(chan struct{})
	wg.Add(2)
	go processSuccessAndError(wg, p)
	go publicMessages(*num, wg, p, donechan)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	close(donechan)
	p.Close()
	wg.Wait()
}
func processSuccessAndError(wg *sync.WaitGroup, p sarama.AsyncProducer) {
	defer wg.Done()
	for {
		select {
		case e, more := <-p.Errors():
			if !more {
				return
			}
			fmt.Printf("Oops! we fail to publish message to kafka broker,err:%s\n", e.Err)
			// TODO: save the message to somewhere for retry later
			// continue
		case s := <-p.Successes():
			key, err := s.Key.Encode()
			if nil != err {
				fmt.Printf("fail to encode key,err:%s", err)
			}
			content, err := s.Value.Encode()
			if nil != err {
				fmt.Printf("fail to encode message value,err:%s\n", err)
			}
			fmt.Printf("we successfully publish msg with key:%s,msg:%s,partition:%d,offset:%d\n", string(key), string(content), s.Partition, s.Offset)
		}
	}
}
func publicMessages(count int, wg *sync.WaitGroup, p sarama.AsyncProducer, donechan chan struct{}) {
	defer wg.Done()
	for i := 0; i < count; i++ {
		pmsg := &sarama.ProducerMessage{
			Topic: *topic,
			Key:   sarama.StringEncoder(fmt.Sprintf("key-%d", i)),
			Value: sarama.StringEncoder(fmt.Sprintf("content-%d", i)),
		}

		select {
		case <-donechan:
			return
		case p.Input() <- pmsg:
		}
	}
}
