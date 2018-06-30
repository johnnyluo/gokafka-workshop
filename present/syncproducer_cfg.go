cfg := sarama.NewConfig()
cfg.ClientID = "my-kafka-producer"
// The following three settings are quite important for Sync Producer
cfg.Producer.Return.Successes = true
cfg.Producer.Return.Errors = true
cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Version = sarama.V1_1_0_0
	c, err := sarama.NewClient([]string{"localhost:9092"}, cfg)
	if nil != err {
		panic(err)
	}
	p, err := sarama.NewSyncProducerFromClient(c)
	if nil != err {
		panic(err)
	}