package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
)

var shards = []string{
	"http://localhost:8081",
	"http://localhost:8082",
	"http://localhost:8083",
}

func main() {
	fmt.Println("Coordinator running on http://localhost:8080")

	http.HandleFunc("/insert", insertHandler)
	http.HandleFunc("/search", searchHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Mapping an ID to a shard index
func getShardIndex(id string) int {
	var sum int
	for _, ch := range id {
		sum += int(ch)
	}
	return sum % len(shards)
}

func insertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//Read
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}

	//Extract ID from JSON
	var payload struct {
		ID string `json:"id"`
	}
	json.Unmarshal(body, &payload)

	if payload.ID == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}

	//Choose shard using hash function
	shardIndex := getShardIndex(payload.ID)
	shardURL := shards[shardIndex] + "/insert"

	//Forward request to shard
	resp, err := http.Post(shardURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, "Error contacting shard", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read shard response
	shardResp, _ := io.ReadAll(resp.Body)

	// Return shard response
	w.Header().Set("Content-Type", "application/json")
	w.Write(shardResp)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}

	// Collect results from all shards
	var allResults []map[string]interface{}

	for _, shard := range shards {
		resp, err := http.Post(shard+"/search", "application/json", bytes.NewBuffer(body))
		if err != nil {
			http.Error(w, "Error contacting shard: "+shard, http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		shardResp, _ := io.ReadAll(resp.Body)

		var shardResults []map[string]interface{}
		json.Unmarshal(shardResp, &shardResults)

		allResults = append(allResults, shardResults...)
	}

	// Sort results by distance
	sort.Slice(allResults, func(i, j int) bool {
		return allResults[i]["distance"].(float64) < allResults[j]["distance"].(float64)
	})

	// Extract top_k from request
	var reqBody struct {
		TopK int `json:"top_k"`
	}
	json.Unmarshal(body, &reqBody)
	topK := reqBody.TopK
	if topK <= 0 || topK > len(allResults) {
		topK = len(allResults)
	}

	finalResults := allResults[:topK]

	// Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(finalResults)
}
