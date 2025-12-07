package router

import (
	"go-webserver/internal/handlers"
	"net/http"
)

func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/hello", handlers.HelloHandler)
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/info", handlers.InfoHandler)
	mux.HandleFunc("/echo", handlers.EchoHandler)
	mux.HandleFunc("/insert", handlers.InsertHandler)
	mux.HandleFunc("/search", handlers.SearchHandler)

	return mux
}
