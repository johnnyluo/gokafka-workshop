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

	// Setup your client id appropriately , otherwise the library will generate a uuid for you
	cfg.ClientID = "my-kafka-producer"

	// If you care about whether your messages had been published successfully,then please set the following two to true
	// trust me , most of the time , you do care
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	// wait for the in sync replica to ack, strong consistency , but slow.
	// This is also related to min.insync.replicas
	// balance between high available / high consistency
	// cluster level or topic level
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	// How frequent to flush the messages to broker. best effort which ever comes first.
	// the way to tune for producing performance
	cfg.Producer.Flush.Bytes = 65535
	cfg.Producer.Flush.Frequency = 500 * time.Millisecond

	// broker accept a slice, make sure you give multiple broker address, or a DNS
	c, err := sarama.NewClient([]string{"127.0.0.1:9092"}, cfg)
	if nil != err {
		panic(err)
	}
	p, err := sarama.NewAsyncProducerFromClient(c)
	if nil != err {
		panic(err)
	}
	defer func() {
		// make sure you do closet the async producer, it is important , otherwise you might risk or losing message
		if err := p.Close(); nil != err {
			fmt.Printf("error while closing producer:%s", err)
		}
		if err := c.Close(); nil != err {
			fmt.Printf("error while closing client:%s", err)
		}
	}()
	wg := &sync.WaitGroup{}
	donechan := make(chan struct{})
	wg.Add(2)
	// start to process success and error notifications in a seperate go routine
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
		// errors will happen, when broker are going through a leader election, it will return EOF on the client library
		case e, more := <-p.Errors():
			if !more {
				return
			}
			fmt.Printf("Oops! we fail to publish message to kafka broker,err:%s\n", e.Err)
			// TODO: save the message to somewhere for retry later
			// continue
		case s, more := <-p.Successes():
			if !more {
				return
			}
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
func publicMessages(count int, wg *sync.WaitGroup, p sarama.AsyncProducer, done chan struct{}) {
	defer wg.Done()
	for i := 0; i < count; i++ {
		pmsg := &sarama.ProducerMessage{
			Topic: *topic,
			Key:   sarama.StringEncoder(fmt.Sprintf("key-%d", i)),
			Value: sarama.StringEncoder(fmt.Sprintf("content-%d", i)),
		}

		select {
		case <-done:
			return
		case p.Input() <- pmsg:
		}
	}
}
