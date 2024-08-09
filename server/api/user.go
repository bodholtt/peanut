package api

import (
	"encoding/json"
	"net/http"
	"peanutserver/database"
	"peanutserver/types"
	"strconv"
)

// API Route: /user

// HandleUser - Route to the proper function depending on the request method.
func HandleUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleUserGET(w, r)
	case http.MethodPost:
		handleUserPOST(w, r)
	}
}

func handleUserGET(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	user, err := database.GetUser(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(types.APIResponse{
		Body:  user,
		Error: "",
	})
}

// handleUserPOST - Handle creation of a new user. Respond with the user's ID if successful.
func handleUserPOST(w http.ResponseWriter, r *http.Request) {
	user, err := database.CreateUser("test", "pass", 1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(types.APIResponse{
		Body:  user,
		Error: "",
	})
}
