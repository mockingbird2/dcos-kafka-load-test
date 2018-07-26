package main

import (
	"flag"
	"fmt"
	"gopkg.in/Shopify/sarama.v1"
	"os"
	"strings"
)

type inputConfig struct {
	brokers      []string
	topic        string
	msgSize      int
	batchSize    int
	compression  string
	requiredAcks string
	bufferSize   int
	msgRate      uint64
	duration     int
	Workers      struct {
		creators  int
		producers int
	}
}

func ParseInput() *inputConfig {
	config := &inputConfig{}
	brokerString := flag.String("brokers", "localhost:9092", "Comma delimited list of Kafka brokers")
	flag.StringVar(&config.topic, "topic", "topic test", "Kafka topic which messages are send to")
	flag.IntVar(&config.msgSize, "message-size", 300, "Message size (bytes)")
	flag.IntVar(&config.batchSize, "message-batch-size", 500, "Messages per batch")

	flag.StringVar(&config.compression, "compression", "none", "Message compression: none, gzip, snappy")
	flag.StringVar(&config.requiredAcks, "required-acks", "local", "RequiredAcks config: none, local, all")

	flag.Uint64Var(&config.msgRate, "produce-rate", 1000, "Global write rate limit (messages/sec)")
	flag.IntVar(&config.bufferSize, "event-buffer-size", 256, "Overall buffered events in produceer")
	flag.IntVar(&config.duration, "duration", 10, "Duration of test in secs")

	flag.IntVar(&config.Workers.producers, "workers", 1, "Number of workers")
	flag.IntVar(&config.Workers.creators, "creators", 1, "Number of message creators")
	flag.Parse()

	config.brokers = strings.Split(*brokerString, ",")
	fmt.Printf("%+v\n", config)
	return config
}

func KafkaConfig(c inputConfig) *sarama.Config {
	conf := sarama.NewConfig()

	switch c.compression {
	case "gzip":
		conf.Producer.Compression = sarama.CompressionGZIP
	case "snappy":
		conf.Producer.Compression = sarama.CompressionSnappy
	case "none":
		conf.Producer.Compression = sarama.CompressionNone
	default:
		fmt.Printf("Invalid compression option: %s\n", c.compression)
		os.Exit(1)
	}

	switch c.requiredAcks {
	case "none":
		conf.Producer.RequiredAcks = sarama.NoResponse
	case "local":
		conf.Producer.RequiredAcks = sarama.WaitForLocal
	case "all":
		conf.Producer.RequiredAcks = sarama.WaitForAll
	default:
		fmt.Printf("Invalid required-acks option: %s\n", c.requiredAcks)
		os.Exit(1)
	}

	conf.Producer.Return.Successes = true
	conf.Producer.Flush.MaxMessages = c.batchSize
	conf.Producer.MaxMessageBytes = c.msgSize + 50
	return conf
}
