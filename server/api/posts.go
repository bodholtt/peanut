package api

import (
	"encoding/json"
	"log"
	"net/http"
	"peanutserver/database"
	"strconv"
)

// HandlePost - Route to the proper function depending on the request method.
func HandlePost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handlePostGET(w, r)
	case http.MethodPost:
		handlePostPOST(w, r)
	}
}

// handlePostGET - handle the delivery of a Post object through the API, for displaying on the client.
func handlePostGET(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Path[len("/post/"):])
	if err != nil {
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")

	post, err := database.GetPost(id)

	json.NewEncoder(w).Encode(post)
}

// handlePostPOST - handle the creation of a Post object in the database.
func handlePostPOST(w http.ResponseWriter, r *http.Request) {

	id, err := database.CreatePost()
	if err != nil {
		log.Fatal(err)
	}
	// TODO: fix this when client and server go in proxy
	location := "http://localhost:4321/post/" + strconv.Itoa(id)

	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusFound)
}

// handlePostPUT - handle the updating of a Post object in the database.
func handlePostPUT(w http.ResponseWriter, r *http.Request) {
	//	to implement
}

// HandlePosts - handle the delivery of a PostThumbs object according to paginated navigation on the client.
func HandlePosts(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		json.NewEncoder(w).Encode("Invalid method")
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

func HandlePostCount(w http.ResponseWriter, r *http.Request) {
	count, err := database.GetPostCount()
	if err != nil {
		log.Println(err)
	}
	json.NewEncoder(w).Encode(count)
}
