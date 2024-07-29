package api

import (
	"encoding/json"
	"net/http"
	"peanutserver/database"
	"strconv"
)

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

	id, err := strconv.Atoi(r.URL.Path[len("/user/"):])
	if err != nil {
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")

	post, err := database.GetUser(id)

	json.NewEncoder(w).Encode(post)
}

func handleUserPOST(w http.ResponseWriter, r *http.Request) {

}
