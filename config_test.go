package main

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/Shopify/sarama.v1"
	"testing"
)

func testCodec(t *testing.T, name string, codec sarama.CompressionCodec) {
	config := &inputConfig{}
	config.requiredAcks = "none"
	config.compression = name
	kafkaConfig := KafkaConfig(*config)
	compression := kafkaConfig.Producer.Compression
	assert.Equal(t, compression, codec)
}

func TestKafkaCompressionConfig(t *testing.T) {
	mapping := map[string]sarama.CompressionCodec{
		"none":   sarama.CompressionNone,
		"snappy": sarama.CompressionSnappy,
		"gzip":   sarama.CompressionGZIP,
	}
	for k, v := range mapping {
		testCodec(t, k, v)
	}
}

func testAck(t *testing.T, name string, ack sarama.RequiredAcks) {
	config := &inputConfig{}
	config.compression = "none"
	config.requiredAcks = name
	kafkaConfig := KafkaConfig(*config)
	ackType := kafkaConfig.Producer.RequiredAcks
	assert.Equal(t, ackType, ack)
}

func TestKafkaAckConfig(t *testing.T) {
	mapping := map[string]sarama.RequiredAcks{
		"none":  sarama.NoResponse,
		"local": sarama.WaitForLocal,
		"all":   sarama.WaitForAll,
	}
	for k, v := range mapping {
		testAck(t, k, v)
	}
}
