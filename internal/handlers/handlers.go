package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

var vectorDB = make(map[string][]float64) // In-memory vector storage

// HelloHandler responds with a simple message
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from your first Go web server!")
}

// HealthHandler responds with OK
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

// InfoHandler responds with JSON info
func InfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{
		"service": "Basic Go Web Server",
		"status":  "running",
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)

}

// EchoHandler handles POST requests and echoes JSON
func EchoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"status": "success",
		"data":   data,
	}

	json.NewEncoder(w).Encode(response)
}

// Vector Insert endpoint

func InsertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var payload struct {
		ID     string    `json:"id"`
		Vector []float64 `json:"vector"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || payload.ID == "" || len(payload.Vector) == 0 {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return

	}

	vectorDB[payload.ID] = payload.Vector

	json.NewEncoder(w).Encode(map[string]string{
		"status": "inserted",
		"id":     payload.ID,
	})
}

// Euclidean distance
func euclidean(a, b []float64) float64 {
	sum := 0.0
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return math.Sqrt(sum)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var query struct {
		Vector []float64 `json:"vector"`
		TopK   int       `json:"top_k"`
	}

	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil || len(query.Vector) == 0 || query.TopK <= 0 {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Find top-k closest vectors
	type result struct {
		ID       string  `json:"id"`
		Distance float64 `json:"distance"`
	}
	var results []result
	for id, vec := range vectorDB {
		if len(vec) != len(query.Vector) {
			continue
		}
		dist := euclidean(vec, query.Vector)
		results = append(results, result{ID: id, Distance: dist})
	}

	// Simple sort by distance
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].Distance > results[j].Distance {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	if len(results) > query.TopK {
		results = results[:query.TopK]
	}

	json.NewEncoder(w).Encode(results)
}
