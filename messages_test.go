package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestCreatorLifecycleFinishs(t *testing.T) {
	config := &inputConfig{}
	config.Workers.creators = 1
	config.msgSize = 1
	config.batchSize = 1
	m := MessageCreator(*config)
	m.StartCreators()
	time.Sleep(1000 * time.Millisecond)
	messages := m.createdMessages
	m.StopCreators()
	assert.True(t, len(messages) > 0)
}

func TestRandomMessageGeneration(t *testing.T) {
	msgSizee := 1000
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	msg1 := randMsg(msgSizee, generator)
	msg2 := randMsg(msgSizee, generator)
	assert.False(t, bytes.Equal(msg1, msg2))
}

func TestMessagePush(t *testing.T) {
	capacity := 10
	m := &messageCreator{}
	messageChannel := make(chan []byte, capacity)
	m.createdMessages = messageChannel
	testMsg := make([]byte, 10)

	err := m.pushMessage(testMsg)
	assert.True(t, err == nil)

	for i := 1; i <= capacity-1; i++ {
		m.createdMessages <- testMsg
	}

	err = m.pushMessage(testMsg)
	assert.True(t, err != nil)

}
