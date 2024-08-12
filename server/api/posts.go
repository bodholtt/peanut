package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"peanutserver/database"
	"peanutserver/types"
	"strconv"
	"strings"
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

	params, _ := url.ParseQuery(r.URL.RawQuery)

	limit, err := strconv.Atoi(params.Get("limit"))
	if err != nil {
		limit = 50
	}
	if limit > 50 || limit < 1 {
		limit = 50
	}

	offset, err := strconv.Atoi(params.Get("offset"))
	if err != nil {
		offset = 0
	}
	if offset < 0 {
		offset = 0
	}

	var tags []string

	tagsQuery := params.Get("tags")
	if tagsQuery != "" {
		tags = strings.Split(tagsQuery, ",")
	}

	//limit <= 50
	//return at most 50 posts
	//if no limit designated return 50 posts

	posts, err := database.GetPostThumbs(limit, offset, tags)
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
