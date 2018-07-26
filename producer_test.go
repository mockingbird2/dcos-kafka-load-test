package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComputeTickerInterval(t *testing.T) {
	assert.Equal(t, computeTickerInterval(1, 1), int64(1e9))
	assert.Equal(t, computeTickerInterval(100, 2), int64(0.02*1e9))
}

func TestMessagePoll(t *testing.T) {
	k := &kafkaProducer{}
	m := make(chan []byte, 10)
	k.messages = m
	_, err := k.pollMessage()
	assert.True(t, err != nil)

	m <- make([]byte, 10)
	_, err = k.pollMessage()
	assert.True(t, err == nil)
}
