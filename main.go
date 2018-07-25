package main

import (
	"fmt"
	"time"
)

func main() {
	config := ParseInput()
	creator := MessageCreator(*config)
	creator.StartCreators()
	producer := KafkaProducer(*config, creator.MessagePool())
	producer.StartProducers()
	timer := time.NewTimer(time.Duration(config.duration) * time.Second)
	<-timer.C
	producer.StopProducers()
	fmt.Println("Stopped producers")
	creator.StopCreators()
	fmt.Println("Stopped creators")
}
