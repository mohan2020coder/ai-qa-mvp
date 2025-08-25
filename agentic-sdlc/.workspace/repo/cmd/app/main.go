package main

import (
	"fmt"
	"log"
	"net/http"
	"app/internal/app/handlers"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok"))
	})
	http.HandleFunc("/users", handlers.Users)

	fmt.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
