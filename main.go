package main

import (
	"fmt"
	"log"
	"net/http"

	"concurrent-port-scanner/internal/handler"
)

func main() {
	mux := http.NewServeMux()
	scanHandler := handler.NewScanHandler()
	mux.HandleFunc("/scan", scanHandler.HandleScan)
	mux.HandleFunc("/", scanHandler.HandleHome)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
