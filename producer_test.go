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
	m := MessagePool(1)
	k.messages = m
	msg := k.pollMessage()
	assert.True(t, msg != nil)
}

func TestClientAmountCalculation(t *testing.T) {
	assert.Equal(t, computeClientAmount(1, 1), 1)
	assert.Equal(t, computeClientAmount(7, 5), 2)

}
