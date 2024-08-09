package api

import (
	"encoding/json"
	"net/http"
	"peanutserver/database"
	"peanutserver/types"
)

// API route: /postCount

func HandlePostCount(w http.ResponseWriter, r *http.Request) {
	count, err := database.GetPostCount()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: err.Error(),
		})
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.APIResponse{
		Body:  count,
		Error: "",
	})
}
