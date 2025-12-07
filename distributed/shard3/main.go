package main

import (
	"fmt"
	"log"
	"net/http"

	"go-webserver/internal/router"
)

func main() {
	r := router.SetupRouter()

	port := "8083"

	fmt.Println("Shard 3 running on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
