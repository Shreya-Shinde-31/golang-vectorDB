package main

import (
	"fmt"
	"log"
	"net/http"

	"go-webserver/internal/router"
)

func main() {
	r := router.SetupRouter() // Initialize router with all routes

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
