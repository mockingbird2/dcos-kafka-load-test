package main

import (
	"fmt"
	"gopkg.in/Shopify/sarama.v1"
	"sync"
	"time"
)

type kafkaProducer struct {
	config   sarama.Config
	input    inputConfig
	messages <-chan []byte
	client   sarama.Client
	ticker   time.Ticker
	wg       sync.WaitGroup
	stop     chan bool
}

func KafkaProducer(config inputConfig, m <-chan []byte) *kafkaProducer {
	c := KafkaConfig(config)
	interval := computeTickerInterval(config.msgRate, config.Workers.producers)
	client, err := sarama.NewClient(config.brokers, c)
	var wg sync.WaitGroup
	if err != nil {
		fmt.Println("New Client Error")
		fmt.Println(err.Error())
	}
	fmt.Println("Connected to kafka client")
	stop := make(chan bool, config.Workers.creators)
	ticker := time.NewTicker(time.Duration(interval) * time.Nanosecond)
	return &kafkaProducer{*c, config, m, client, *ticker, wg, stop}
}

func (k *kafkaProducer) StartProducers() {
	count := k.input.Workers.producers
	k.wg.Add(count)
	for i := 1; i <= count; i++ {
		go k.producer()
	}
}

func (k *kafkaProducer) StopProducers() {
	k.ticker.Stop()
	for i := 1; i <= k.input.Workers.producers; i++ {
		k.stop <- true
	}
	k.wg.Wait()
}

func (k *kafkaProducer) producer() {
	p, err := sarama.NewSyncProducerFromClient(k.client)
	defer p.Close()
	defer k.wg.Done()
	msgBatch := make([]*sarama.ProducerMessage, 0, k.input.batchSize)
	if err != nil {
		fmt.Println("New Producer Error")
		fmt.Println(err.Error())
		return
	}
	for range k.ticker.C {
		select {
		case m := <-k.messages:
			msgBatch = append(msgBatch, BuildProducerMessage(k.input.topic, m))
			if len(msgBatch) != k.input.batchSize {
				continue
			}
			err := p.SendMessages(msgBatch)
			msgBatch = msgBatch[:0]
			if err != nil {
				fmt.Println("Error while sending")
			}
		default:
		}
		select {
		case <-k.stop:
			fmt.Println("Stopped producer")
			return
		default:
			continue
		}
	}
}

func computeTickerInterval(msgRate uint64, workerCount int) int64 {
	return int64((1.0 / (float64(msgRate) / float64(workerCount))) * 1e9)
}
