package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	r := &Repository{}
	r.Init()

	http.HandleFunc("/stormReports", handleGetStormReports(r))

	fmt.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
