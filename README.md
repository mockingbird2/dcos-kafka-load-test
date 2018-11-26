# dcos-kafka-load-test
[![Build Status](https://travis-ci.org/fapaul/dcos-kafka-load-test.svg?branch=master)](https://travis-ci.org/fapaul/dcos-kafka-load-test)
[![Go Report Card](https://goreportcard.com/badge/github.com/mockingbird2/dcos-kafka-load-test)](https://goreportcard.com/report/github.com/mockingbird2/dcos-kafka-load-test)

## Usage

```
Usage of ./dcos-kafka-load-test:
  -brokers string
    	Comma delimited list of Kafka brokers (default "localhost:9092")
  -compression string
    	Message compression: none, gzip, snappy (default "none")
  -creators int
    	Number of message creators (default 1)
  -duration int
    	Duration of test in secs (default 10)
  -event-buffer-size int
    	Overall buffered events in produceer (default 256)
  -message-batch-size int
    	Messages per batch (default 500)
  -message-size int
    	Message size (bytes) (default 300)
  -produce-rate uint
    	Global write rate limit (messages/sec) (default 1000)
  -required-acks string
    	RequiredAcks config: none, local, all (default "local")
  -topic string
    	Kafka topic which messages are send to (default "topic test")
  -workers int
    	Number of workers (default 1)
```

## Deployment

First start your dcos cluster and install Apache Kafka.

```
deploy.py

Usage:
    deploy.py [options]

Arguments:

Options:
    --dcos-username <user>        DC/OS username
    --dcos-password <pass>        DC/OS password
    --script-cpus <n>             number of CPUs to use to run the script [default: 4]
    --script-mem <n>              amount of memory (mb) to use to run the script [default: 4096]
    --brokers <s>                 Comma delimited list of Kafka brokers [default: localhost:9092]
    --compression <s>             Message compression: none, gzip, snappy [default: none]
    --creators <n>                Number of message creators [default: 1]
    --duration <n>                Duration of test in secs [default: 10]
    --event-buffer-size <n>       Overall buffered events in produceer [default: 256]
    --message-batch-size <n>      Messages per batch [default: 500]
    --message-size <n>            Message size (bytes) [default: 300]
    --produce-rate <n>            Global write rate limit (messages/sec) [default: 1000]
    --required-acks <s>           RequiredAcks config: none, local, all [default: local]
    --topic <s>                   Kafka topic which messages are send to [default: topic_test]
    --workers <n>                 Number of workers [default: 1]
```
