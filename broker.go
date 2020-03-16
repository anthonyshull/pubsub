package main

import "sync"

type Broker struct {
	mu          sync.Mutex
	subscribers map[string]chan []byte
}

func (broker *Broker) Add(key string) {
	broker.mu.Lock()
	defer broker.mu.Unlock()
	broker.subscribers[key] = make(chan []byte)
}

func (broker *Broker) Remove(key string) {
	broker.mu.Lock()
	defer broker.mu.Unlock()
	delete(broker.subscribers, key)
}

func (broker *Broker) Publish(message []byte) {
	broker.mu.Lock()
	defer broker.mu.Unlock()
	for _, c := range broker.subscribers {
		c <- message
	}
}
