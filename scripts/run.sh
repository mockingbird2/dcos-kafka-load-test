#!/bin/sh

brokers=$BROKERS
topic=$TOPIC
mSize=$MESSAGE_SIZE
bSize=$BATCH_SIZE
compr=$COMPRESSION
acks=$ACKS
mRate=$MESSAGE_RATE
eBufferSize=$BUFFER_SIZE
duration=$DURATION
workers=$WORKERS
creators=$CREATORS

/go/src/dcos-kafka-load-test/dcos-kafka-load-test \
    -brokers $brokers \
    -topic $topic \
    -message-size $mSize \
    -message-batch-size $bSize \
    -compression $compr \
    -required-acks $acks \
    -produce-rate $mRate \
    -event-buffer-size $eBufferSize \
    -duration $duration \
    -workers $workers \
    -creators $creators
