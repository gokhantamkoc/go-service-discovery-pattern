package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type Service struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Registry struct {
	services map[string][]Service
	mu       sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		services: make(map[string][]Service),
	}
}

func (r *Registry) registerService(w http.ResponseWriter, req *http.Request) {
	var service Service
	if err := json.NewDecoder(req.Body).Decode(&service); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	r.mu.Lock()
	r.services[service.Name] = append(r.services[service.Name], service)
	r.mu.Unlock()
	w.WriteHeader(http.StatusOK)
}

func (r *Registry) listServices(w http.ResponseWriter, req *http.Request) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	services := make([]Service, 0)
	for _, s := range r.services {
		services = append(services, s...)
	}
	json.NewEncoder(w).Encode(services)
}

func main() {
	registry := NewRegistry()

	http.HandleFunc("/register", registry.registerService)
	http.HandleFunc("/services", registry.listServices)

	server := &http.Server{
		Addr:           ":8080",
		Handler:        http.DefaultServeMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Service Registry is running at :8080")
	log.Fatal(server.ListenAndServe())
}
