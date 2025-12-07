package main

import (
	"fmt"
	"go-webserver/internal/router"
	"log"
	"net/http"
)

func main() {
	r := router.SetupRouter()

	port := "8081"

	fmt.Println("Shard 1 running on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
