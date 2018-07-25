package main

import (
	"fmt"
	"gopkg.in/Shopify/sarama.v1"
	"math/rand"
	"sync"
	"time"
)

type messageCreatorStrategy func([]byte)

type messageCreator struct {
	config          inputConfig
	createdMessages chan<- []byte
	wg              sync.WaitGroup
	stop            chan bool
	messagePool     <-chan []byte
}

func MessageCreator(config inputConfig) *messageCreator {
	var wg sync.WaitGroup
	messages := make(chan []byte, config.batchSize*100)
	stop := make(chan bool, config.Workers.creators)
	return &messageCreator{config, messages, wg, stop, messages}
}

func (m *messageCreator) StopCreators() {
	config := m.config
	for i := 1; i <= config.Workers.creators; i++ {
		m.stop <- true
	}
	m.wg.Wait()
}

func (m *messageCreator) StartCreators() {
	config := m.config
	m.wg.Add(config.Workers.creators)
	for i := 1; i <= config.Workers.creators; i++ {
		go m.creator()
	}
}

func (m *messageCreator) creator() {
	config := m.config
	defer m.wg.Done()
	msg := createMessage(config.msgSize, randMsg)
	for {
		select {
		case m.createdMessages <- msg:
		default:
		}
		select {
		case <-m.stop:
			fmt.Println("Stopped message creator")
			return
		default:
		}
	}
}

func createMessage(size int, fn messageCreatorStrategy) []byte {
	message := make([]byte, size)
	fn(message)
	return message
}

func BuildProducerMessage(topic string, msgData []byte) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{Topic: topic, Value: sarama.ByteEncoder(msgData)}
	return msg

}

func randMsg(m []byte) {
	chars := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$^&*(){}][:<>.")
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	for i := range m {
		m[i] = chars[generator.Intn(len(chars))]
	}
}

func (m *messageCreator) MessagePool() <-chan []byte {
	return m.messagePool
}
