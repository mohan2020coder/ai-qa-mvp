package handlers

import (
	"encoding/json"
	"net/http"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func Users(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode([]User{{ID:1, Name:"Ada"}})
}
