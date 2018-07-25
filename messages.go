package main

import (
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
}

func MessageCreator(config inputConfig, messages chan<- []byte) *messageCreator {
	var wg sync.WaitGroup
	stop := make(chan bool, config.Workers.creators)
	return &messageCreator{config, messages, wg, stop}
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
	msgData := make([]byte, config.msgSize)
	createMessage(msgData, randMsg)
	for {
		select {
		case m.createdMessages <- msgData:
		default:
		}
		select {
		case <-m.stop:
			return
		default:
		}
	}
}

func createMessage(message []byte, fn messageCreatorStrategy) {
	fn(message)
}

func randMsg(m []byte) {
	chars := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$^&*(){}][:<>.")
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	for i := range m {
		m[i] = chars[generator.Intn(len(chars))]
	}
}
