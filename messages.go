package main

import (
	"errors"
	"gopkg.in/Shopify/sarama.v1"
	"math/rand"
	"sync"
	"time"
)

type messageCreator struct {
	config          inputConfig
	createdMessages chan<- []byte
	wg              *sync.WaitGroup
	stop            chan bool
	messagePool     <-chan []byte
}

func MessageCreator(config inputConfig) *messageCreator {
	var wg sync.WaitGroup
	messages := make(chan []byte, config.batchSize*100)
	stop := make(chan bool, config.Workers.creators)
	return &messageCreator{config, messages, &wg, stop, messages}
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
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	msg := randMsg(config.msgSize, generator)
	for {
		m.pushMessage(msg)
		select {
		case <-m.stop:
			return
		default:
		}
	}
}

func BuildProducerMessage(topic string, msgData []byte) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{Topic: topic, Value: sarama.ByteEncoder(msgData)}
	return msg

}

func (m *messageCreator) pushMessage(msg []byte) error {
	select {
	case m.createdMessages <- msg:
		return nil
	default:
		return errors.New("Message was not pushed, discarded")
	}

}

func randMsg(size int, generator *rand.Rand) []byte {
	m := make([]byte, size)
	chars := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$^&*(){}][:<>.")
	for i := range m {
		m[i] = chars[generator.Intn(len(chars))]
	}
	return m
}

func (m *messageCreator) MessagePool() <-chan []byte {
	return m.messagePool
}
