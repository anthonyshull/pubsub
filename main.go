package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	broker := Broker{subscribers: make(map[string]chan []byte)}

	r := chi.NewRouter()
	r.Post("/publish", broker.Publish)
	r.Get("/subscribe", broker.Subscribe)

	log.Fatal(http.ListenAndServe(":9999", r))
}
