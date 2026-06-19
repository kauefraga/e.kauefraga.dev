package main

import (
	"fmt"
	"net/http"

	"github.com/kauefraga/e.kauefraga.dev/internal/shortener"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/links", shortener.CreateLink)
	mux.HandleFunc("GET /{code}", shortener.RetrieveLink)

	fmt.Println("Server is listening on http://localhost:3333")
	http.ListenAndServe(":3333", mux)
}
