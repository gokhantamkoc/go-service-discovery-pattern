package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Service struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func main() {
	service := Service{
		Name: "example-service",
		URL:  "http://localhost:8081",
	}

	data, err := json.Marshal(service)
	if err != nil {
		log.Fatalf("Error marshaling service: %v", err)
	}

	for {
		resp, err := http.Post("http://registry:8080/register", "application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Printf("Error registering service: %v", err)
		} else {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				log.Println("Service registered successfully")
				break
			}
		}
		time.Sleep(5 * time.Second) // Retry every 5 seconds
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request")
		w.Write([]byte("Hello from example-service"))
	})

	log.Println("Example service is running at :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
