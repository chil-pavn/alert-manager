package main

import (
	"log"
	"net/http"

	"github.com/chil-pavn/alert-manager/receiver"
)

func main() {
	http.HandleFunc("/webhook", receiver.HandleWebhook)
	port := ":5000"
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
