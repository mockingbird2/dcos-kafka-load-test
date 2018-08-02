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
	messages := m.Pool()
	m.StopCreators()
	assert.True(t, messages.Get() != nil)
}

func TestRandomMessageGeneration(t *testing.T) {
	msgSize := 1000
	msg1 := make([]byte, msgSize)
	msg2 := make([]byte, msgSize)
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	randMsg(msg1, generator)
	randMsg(msg2, generator)
	assert.False(t, bytes.Equal(msg1, msg2))
}
