package api

import (
	"encoding/json"
	"net/http"
	"peanutserver/database"
	"peanutserver/types"
	"strconv"
)

// API route: /post/{id}/tags

// HandlePostTags - Route to the proper function depending on the request method for the /post/{id}/tags route
func HandlePostTags(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, OPTIONS")

	switch r.Method {
	case http.MethodGet:
		handlePostTagsGET(w, r)
		//case http.MethodPut:
		//	handlePostTagsPUT(w, r)
	}
}

// handlePostTagsGET - get the tags of a post by the post ID
func handlePostTagsGET(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: err.Error(),
		})
		return
	}

	tags, err := database.GetTagsByPostID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: err.Error(),
		})
	}

	json.NewEncoder(w).Encode(types.APIResponse{
		Body:  tags,
		Error: "",
	})
}

// handlePostTagsPUT - handle the updating of a Post's tags.
//func handlePostTagsPUT(w http.ResponseWriter, r *http.Request) {
//
//	w.Header().Set("Access-Control-Allow-Origin", pcfg.Cfg.Client.Host)
//
//	_, err := strconv.Atoi(r.PathValue("id"))
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		json.NewEncoder(w).Encode(Response{
//			Body:  nil,
//			Error: err.Error(),
//		})
//		return
//	}
//
//	// Accept a string with a list of tag names. Get the IDs of each tag - if the tag doesn't exist, create it
//	// Update post_tags with the id of the post and the ids of the new tags
//	err = r.ParseForm()
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		json.NewEncoder(w).Encode(Response{
//			Body:  nil,
//			Error: err.Error(),
//		})
//		return
//	}
//	log.Println(r.PostFormValue("tags"))
//	// TODO: THIS DOESNT FUCKING WORK!!!!!!!
//
//	json.NewEncoder(w).Encode("Not implemented")
//}
