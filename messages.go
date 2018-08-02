package main

import (
	"fmt"
	"gopkg.in/Shopify/sarama.v1"
	"math/rand"
	"sync"
	"time"
)

type messageCreator struct {
	config      inputConfig
	wg          *sync.WaitGroup
	stop        chan bool
	messagePool *sync.Pool
}

func MessageCreator(config inputConfig) *messageCreator {
	var wg sync.WaitGroup
	messages := MessagePool(config.msgSize)
	stop := make(chan bool, config.Workers.creators)
	return &messageCreator{config, &wg, stop, messages}
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
	defer m.wg.Done()
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	msg := make([]byte, m.config.msgSize)
	randMsg(msg, generator)
	for {
		m.pushMessage(msg)
		select {
		case <-m.stop:
			fmt.Println("Stopped creator")
			return
		default:
		}
	}
}

func BuildProducerMessage(topic string, msgData []byte) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{Topic: topic, Value: sarama.ByteEncoder(msgData)}
	return msg

}

func (m *messageCreator) pushMessage(msg []byte) {
	m.Pool().Put(msg)
}

func randMsg(m []byte, generator *rand.Rand) {
	chars := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$^&*(){}][:<>.")
	for i := range m {
		m[i] = chars[generator.Intn(len(chars))]
	}
}

func (m *messageCreator) Pool() *sync.Pool {
	return m.messagePool
}
