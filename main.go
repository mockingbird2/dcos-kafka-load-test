package main

import "time"

func main() {
	config := ParseInput()
	creator := MessageCreator(*config)
	creator.StartCreators()
	timer := time.NewTimer(time.Duration(config.duration) * time.Second)
	<-timer.C
	creator.StopCreators()
}
