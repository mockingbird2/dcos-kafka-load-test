package main

import (
	"fmt"
	"gopkg.in/Shopify/sarama.v1"
	"time"
)

type kafkaProducer struct {
	config   sarama.Config
	input    inputConfig
	messages <-chan *sarama.ProducerMessage
	client   sarama.Client
	ticker   time.Ticker
}

func KafkaProducer(config inputConfig, m <-chan *sarama.ProducerMessage) *kafkaProducer {
	c := KafkaConfig(config)
	interval := computeTickerInterval(config.msgRate, config.Workers.producers)
	client, err := sarama.NewClient(config.brokers, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return &kafkaProducer{*c, config, m, client, *time.NewTicker(time.Duration(interval) * time.Nanosecond)}
}

func (k *kafkaProducer) StartProducers() {
}

func (k *kafkaProducer) producer() {
	p, err := sarama.NewSyncProducerFromClient(k.client)
	defer p.Close()
	msgBatch := make([]*sarama.ProducerMessage, 0, k.input.batchSize)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for {
		select {
		case m := <-k.messages:
			msgBatch = append(msgBatch, m)
			if len(msgBatch) != k.input.batchSize {
				continue
			}
			err = p.SendMessages(msgBatch)
			if err != nil {
				fmt.Println("Error while sending")
			}
		default:
		}
	}
}

func computeTickerInterval(msgRate uint64, workerCount int) int64 {
	return int64((1.0 / (float64(msgRate) / float64(workerCount))) * 1e9)
}
