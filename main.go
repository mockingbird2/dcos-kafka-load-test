package main

import (
	"time"
)

func main() {
	InitLoggers()
	config := ParseInput()
	creator := MessageCreator(*config)
	creator.StartCreators()
	producer := KafkaProducer(*config, creator.MessagePool())
	producer.StartProducers()
	timer := time.NewTimer(time.Duration(config.duration) * time.Second)
	<-timer.C
	producer.StopProducers()
	LogInfo("Stopped producers")
	creator.StopCreators()
	LogInfo("Stopped creators")
}
