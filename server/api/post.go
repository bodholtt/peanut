package api

// API route: /post

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"peanutserver/database"
	"strconv"
)

// HandlePost - Route to the proper function depending on the request method for the /post route
func HandlePost(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	switch r.Method {
	case http.MethodGet:
		handlePostGET(w, r)
	case http.MethodPost:
		handlePostPOST(w, r)
	case http.MethodPut:
		handlePostPUT(w, r)
	case http.MethodOptions:
		handlePostOPTIONS(w, r)
	case http.MethodDelete:
		handlePostDELETE(w, r)
	}
}

func handlePostOPTIONS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

// handlePostGET - handle the delivery of a Post object through the API, for displaying on the client.
func handlePostGET(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Path[len("/post/"):])
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	post, err := database.GetPost(id)
	if err != nil {
		http.NotFound(w, r)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	prev, next := database.GetNextAndPreviousPostIDs(id, "")

	post.Previous = strconv.Itoa(prev)
	post.Next = strconv.Itoa(next)

	json.NewEncoder(w).Encode(post)
}

// handlePostPOST - handle the creation of a Post object in the database.
func handlePostPOST(w http.ResponseWriter, r *http.Request) {

	// TODO: alter this so that only the client can make the request
	w.Header().Set("Access-Control-Allow-Origin", "*")

	imageFile, header, err := r.FormFile("image")
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}

	id, imagePath, err := database.CreatePost(filepath.Ext(header.Filename))
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}

	err = uploadImage(imageFile, imagePath)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}

	json.NewEncoder(w).Encode(id)
}

// handlePostPUT - handle the updating of a Post object in the database.
func handlePostPUT(w http.ResponseWriter, r *http.Request) {
	//	to implement
	json.NewEncoder(w).Encode("Not implemented")
}

// handlePostDELETE - handle the deletion of a Post object in the database.
func handlePostDELETE(w http.ResponseWriter, r *http.Request) {

	// TODO: alter this so that only the client can make the request
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id, err := strconv.Atoi(r.URL.Path[len("/post/"):])
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}

	err = database.DeletePost(id)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}

	w.WriteHeader(http.StatusOK)
}

// uploadImage - handle uploading an image to local storage
func uploadImage(fileData multipart.File, filepath string) error {

	localFile, err := os.Create("." + filepath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, fileData)
	if err != nil {
		return err
	}
	return nil
}
