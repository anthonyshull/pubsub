package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// subscribe to messages
func subscribe(broker *Broker) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
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
		broker.Add(key)

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
}


// publish a message
func publish(broker *Broker) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "message body could not be read", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		go broker.Publish(body)
	}
}

func main() {
	broker := &Broker{subscribers: make(map[string] chan[]byte)}

	r := chi.NewRouter()
	r.Post("/publish", publish(broker))
	r.Get("/subscribe", subscribe(broker))

	log.Fatal(http.ListenAndServe(":9999", r))
}
