package main

import (
	"fmt"
	"log"
	"net/http"

	"go-webserver/internal/router"
)

func main() {
	r := router.SetupRouter()

	port := "8082"

	fmt.Println("Shard 2 running on Port:" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
