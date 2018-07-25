package main

import (
	"github.com/stretchr/testify/assert"
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
