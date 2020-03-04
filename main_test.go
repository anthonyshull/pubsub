package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlers(t *testing.T) {

	t.Run("test that we can publish a message", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/publish", strings.NewReader("this is my test data"))
		if err != nil {
			t.Error("request failed")
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(publish)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("wrong status code: got %v want %v", status, http.StatusCreated)
		}
	})

	t.Run("test that we can subscribe to messages", func(t *testing.T) {
		// here i would have to test publishing and subscribing
	})
}
