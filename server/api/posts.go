package api

import (
	"encoding/json"
	"net/http"
	"peanutserver/database"
	"peanutserver/types"
	"strconv"
)

// API route: /posts

// HandlePosts - Route to the proper function depending on the request method for the /posts route
func HandlePosts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handlePostsGET(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

// handlePostsGET - handle the delivery of a PostThumbs object according to paginated navigation on the client.
func handlePostsGET(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 50
	}
	if limit > 50 || limit < 1 {
		limit = 50
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}
	if offset < 0 {
		offset = 0
	}

	//limit <= 50
	//return at most 50 posts
	//if no limit designated return 50 posts

	posts, err := database.GetPostThumbs(limit, offset)
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
		Body:  posts,
		Error: "",
	})
}
