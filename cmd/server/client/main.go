package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Service struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func main() {
	// Step 1: Query the service registry
	resp, err := http.Get("http://localhost:8080/services")
	if err != nil {
		log.Fatalf("Error querying service registry: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Service registry returned non-OK status: %v", resp.Status)
	}

	var services []Service
	if err := json.NewDecoder(resp.Body).Decode(&services); err != nil {
		log.Fatalf("Error decoding service registry response: %v", err)
	}

	if len(services) == 0 {
		log.Fatalf("No services found in the registry")
	}

	// Step 2: Select a service (just use the first one for simplicity)
	selectedService := services[0]

	log.Printf("Selected service: %s\n", selectedService.URL)
	// Step 3: Make a request to the selected service
	serviceResp, err := http.Get(selectedService.URL)
	if err != nil {
		log.Fatalf("Error making request to service: %v", err)
	}
	defer serviceResp.Body.Close()

	if serviceResp.StatusCode != http.StatusOK {
		log.Fatalf("Service returned non-OK status: %v", serviceResp.Status)
	}

	var responseBody []byte
	serviceResp.Body.Read(responseBody)

	fmt.Printf("Response from service (%s): %s\n", selectedService.URL, string(responseBody))
}
