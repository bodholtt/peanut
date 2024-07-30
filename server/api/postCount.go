package api

// API route: /postCount

import (
	"encoding/json"
	"log"
	"net/http"
	"peanutserver/database"
)

func HandlePostCount(w http.ResponseWriter, r *http.Request) {
	count, err := database.GetPostCount()
	if err != nil {
		log.Println(err)
	}
	json.NewEncoder(w).Encode(count)
}
