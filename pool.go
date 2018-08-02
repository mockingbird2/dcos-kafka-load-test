package main

import (
	"math/rand"
	"sync"
	"time"
)

func MessagePool(size int) *sync.Pool {
	messagePool := sync.Pool{
		New: func() interface{} {
			msg := make([]byte, size)
			source := rand.NewSource(time.Now().UnixNano())
			generator := rand.New(source)
			randMsg(msg, generator)
			return msg
		},
	}
	return &messagePool
}
