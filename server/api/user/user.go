package user

import (
	"encoding/json"
	"net/http"
	"peanutserver/database"
	"peanutserver/pcfg"
	"peanutserver/types"
	"strconv"
)

// API Route: /user

func HandleUserOPTIONS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", pcfg.Cfg.Client.Host)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
}

// HandleUserGET - Retrieve a user's information from the database and return it in the response body
func HandleUserGET(w http.ResponseWriter, r *http.Request) {

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

// HandleUserPOST - Handle creation of a new user. Respond with the user's ID if successful.
func HandleUserPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", pcfg.Cfg.Client.Host)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	user := &types.User{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: err.Error(),
		})
		return
	}

	// a user cannot be rank 0 (anonymous)
	if user.Rank == 0 {
		user.Rank = pcfg.Perms.DefaultRank
	}

	hashedPassword := hashPassword(user.Password)
	userID, err := database.CreateUser(user.Username, hashedPassword, user.Rank)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(types.APIResponse{
		Body:  userID,
		Error: "",
	})
}
