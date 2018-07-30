package main

import (
	"errors"
	"fmt"
	"gopkg.in/Shopify/sarama.v1"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

const ProducerPerWorker float64 = 5.0

type kafkaProducer struct {
	input    inputConfig
	messages <-chan []byte
	clients  []sarama.Client
	wg       *sync.WaitGroup
	stop     chan bool
	metrics  *producerMetrics
	interval int64
}

type producerMetrics struct {
	sentBatches uint64
	errors      uint64
}

func KafkaProducer(config inputConfig, m <-chan []byte) *kafkaProducer {
	c := KafkaConfig(config)
	interval := computeTickerInterval(config.msgRate, config.Workers.producers)
	clientCount := computeClientAmount(float64(config.Workers.producers), ProducerPerWorker)
	clients, err := createClients(clientCount, config.brokers, c)
	var wg sync.WaitGroup
	if err != nil {
		fmt.Println("New Client Error")
		fmt.Println(err.Error())
	}
	fmt.Println("Connected to kafka client")
	stop := make(chan bool, config.Workers.creators)
	return &kafkaProducer{config, m, clients, &wg, stop, initMetrics(), interval}
}

func initMetrics() *producerMetrics {
	return &producerMetrics{0, 0}
}

func (k *kafkaProducer) StartProducers() {
	count := k.input.Workers.producers
	k.wg.Add(count)
	for i := 0; i < count; i++ {
		client := k.clients[i/int(ProducerPerWorker)]
		go k.producer(client)
	}
}

func (k *kafkaProducer) StopProducers() {
	for i := 1; i <= k.input.Workers.producers; i++ {
		k.stop <- true
	}
	k.wg.Wait()
	batchCount := atomic.LoadUint64(&k.metrics.sentBatches)
	sendingErrors := atomic.LoadUint64(&k.metrics.errors)
	fmt.Println("Sent batches: ", batchCount)
	fmt.Println("Errors while sending: ", sendingErrors)
}

func (k *kafkaProducer) producer(client sarama.Client) {
	p, err := sarama.NewSyncProducerFromClient(client)
	defer k.wg.Done()
	defer p.Close()
	if err != nil {
		fmt.Println("New Producer Error")
		fmt.Println(err.Error())
		return
	}
	k.startSchedule(p)
}

func createClients(amount int, brokers []string, config *sarama.Config) ([]sarama.Client, error) {
	clients := make([]sarama.Client, amount, amount)
	for i := 0; i < amount; i++ {
		client, err := sarama.NewClient(brokers, config)
		if err != nil {
			return clients, err
		}
		clients[i] = client
	}
	return clients, nil
}

func computeClientAmount(workers float64, writers float64) int {
	return int(math.Ceil(workers / writers))
}

func (k *kafkaProducer) startSchedule(p sarama.SyncProducer) {
	ticker := time.NewTicker(time.Duration(k.interval) * time.Nanosecond)
	msgBatch := make([]*sarama.ProducerMessage, 0, k.input.batchSize)
	for range ticker.C {
		m, err := k.pollMessage()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			msgBatch = append(msgBatch, BuildProducerMessage(k.input.topic, m))
			if len(msgBatch) != k.input.batchSize {
				continue
			}
			atomic.AddUint64(&k.metrics.sentBatches, 1)
			err = p.SendMessages(msgBatch)
			msgBatch = msgBatch[:0]
			if err != nil {
				atomic.AddUint64(&k.metrics.errors, 1)
				fmt.Println("Error while sending")
			}
		}
		select {
		case <-k.stop:
			fmt.Println("Stopped producer")
			ticker.Stop()
			return
		default:
			continue
		}
	}
}

func (k *kafkaProducer) pollMessage() ([]byte, error) {
	select {
	case m := <-k.messages:
		return m, nil
	default:
		return nil, errors.New("Queue not pollable")
	}
}

func computeTickerInterval(msgRate uint64, workerCount int) int64 {
	return int64((1.0 / (float64(msgRate) / float64(workerCount))) * 1e9)
}
