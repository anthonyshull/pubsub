package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type Broker struct {
	mu          sync.Mutex
	subscribers map[string]chan []byte
}

func (broker *Broker) add(key string) {
	broker.mu.Lock()
	defer broker.mu.Unlock()
	broker.subscribers[key] = make(chan []byte)
}

func (broker *Broker) remove(key string) {
	broker.mu.Lock()
	defer broker.mu.Unlock()
	delete(broker.subscribers, key)
}

func (broker *Broker) publish(message []byte) {
	broker.mu.Lock()
	defer broker.mu.Unlock()
	for _, c := range broker.subscribers {
		c <- message
	}
}

func (broker *Broker) Subscribe(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "server does not support streaming", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	key := r.RemoteAddr
	broker.add(key)
	defer broker.remove(key)

	for {
		select {
		case message := <-broker.subscribers[key]:
			fmt.Fprintf(w, "%s\n", message)
			flusher.Flush()
		case <-ctx.Done():
			return
		}
	}
}

func (broker *Broker) Publish(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "message body could not be read", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	go broker.publish(body)
}