package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

var broker = make(map[string]chan []byte)

// subscribe to messages
func subscribe(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "client does not support streaming", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	key := uuid.New().String()
	broker[key] = make(chan []byte)
	defer delete(broker, key)

	for {
		fmt.Fprintf(w, "%s\n", <-broker[key])
		flusher.Flush()
	}
}

// publish a message
func publish(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err == nil {
		for _, c := range broker {
			c <- body
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func main() {
	r := chi.NewRouter()
	r.Post("/publish", publish)
	r.Get("/subscribe", subscribe)
	log.Fatal(http.ListenAndServe(":9999", r))
}
