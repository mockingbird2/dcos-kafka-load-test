package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	createdMessages = make(chan []byte, 100)
)

type messageCreator func([]byte)

func StopCreators(stop chan<- bool, wg *sync.WaitGroup) {
	for i := 1; i <= config.Workers.creators; i++ {
		stop <- true
	}
	fmt.Println("Queued messages %d", len(createdMessages))
	wg.Wait()
}

func StartCreators(stop <-chan bool, wg *sync.WaitGroup) {
	wg.Add(config.Workers.creators)
	for i := 1; i <= config.Workers.creators; i++ {
		go creator(stop, wg)
	}
}

func creator(stop <-chan bool, wg *sync.WaitGroup) {
	//defer wg.Done()
	msgData := make([]byte, config.msgSize)
	for {
		createdMessages <- msgData
		select {
		case <-stop:
			fmt.Println("Stopping Creator")
			wg.Done()
		default:
			fmt.Println("Running Creator")
		}
		//time.Sleep(250 * time.Millisecond)
	}
}

func createMessage(message []byte, fn messageCreator) {
	fn(message)
}

func randMsg(m []byte) {
	chars := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$^&*(){}][:<>.")
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	for i := range m {
		m[i] = chars[generator.Intn(len(chars))]
	}
}
