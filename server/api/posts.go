package api

// API route: /posts

import (
	"encoding/json"
	"net/http"
	"peanutserver/database"
	"strconv"
)

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

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	posts, _ := database.GetPostThumbs(limit, offset)
	json.NewEncoder(w).Encode(posts)
}
