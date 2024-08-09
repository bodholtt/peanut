package api

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"peanutserver/auth"
	"peanutserver/database"
	"peanutserver/files"
	"peanutserver/pcfg"
	"peanutserver/types"
	"strconv"
	"strings"
)

// API route: /post

// HandlePostOPTIONS - OPTIONS for /post
func HandlePostOPTIONS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", pcfg.Cfg.Client.Host)
	w.Header().Set("Access-Control-Allow-Methods", "GET, DELETE, PUT, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// HandlePostGET - handle the delivery of a Post object through the API
func HandlePostGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Authorization")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: err.Error(),
		})
		return
	}

	post, err := database.GetPost(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	prev, next := database.GetNextAndPreviousPostIDs(id, "")

	post.Previous = strconv.Itoa(prev)
	post.Next = strconv.Itoa(next)

	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(post)
}

// HandlePostPOST - handle the creation of a Post object in the database.
func HandlePostPOST(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", pcfg.Cfg.Client.Host)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Check headers and verify tokens ---------------------------------

	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: "content-Type must be multipart/form-data",
		})
		return
	}

	userID, err := auth.GetUserIDFromAuthHeader(r)
	if err != nil && pcfg.Cfg.Permissions.CreatePosts != 0 {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  "",
			Error: "insufficient permission to create a new post",
		})
		return
	}

	// Parse form data --------------------------------------------------

	imageFile, header, err := r.FormFile("image")
	if err != nil {
		log.Println("failed creating post at r.FormFile", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  "",
			Error: err.Error(),
		})
		return
	}

	// TODO: Implement md5 hash of the image and store with post in db

	//var buf1, buf2 bytes.Buffer
	//writer := io.MultiWriter(&buf1, &buf2)

	//hash := md5.New()

	// Create the post ----------------------------------------------------------------

	id, imagePath, err := database.CreatePost(filepath.Ext(header.Filename), userID)
	if err != nil {
		log.Println("failed creating post at database.CreatePost", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: err.Error(),
		})
		return
	}

	err = files.UploadImage(imageFile, imagePath)
	if err != nil {
		log.Println("failed creating post at files.UploadImage", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(types.APIResponse{
		Body:  id,
		Error: "",
	})
}

// putData - struct for accepting data from a PUT request, handled by HandlePostPUT.
type putData struct {
	Tags   string `json:"tags"`
	Source string `json:"source"`
}

// HandlePostPUT - handle the updating of a Post object in the database.
func HandlePostPUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", pcfg.Cfg.Client.Host)

	// Really shitty way to differentiate between requests to allow
	// for omission in order to not update the value
	// and for a value to be updated as empty
	data := &putData{
		Tags:   "UNGUESSABLE_DEFAULT_VALUE",
		Source: "UNGUESSABLE_DEFAULT_VALUE",
	}

	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: err.Error(),
		})
		return
	}

	log.Println(data.Tags, " , ", data.Source)

	if data.Tags != "UNGUESSABLE_DEFAULT_VALUE" {
		log.Println("Call database function to update tags")
	}
	if data.Source != "UNGUESSABLE_DEFAULT_VALUE" {
		log.Println("Call database function to update source")
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Not implemented")
}

// HandlePostDELETE - handle the deletion of a Post object in the database.
func HandlePostDELETE(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", pcfg.Cfg.Client.Host)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: err.Error(),
		})
		return
	}

	filename, err := database.DeletePost(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: err.Error(),
		})
		return
	}

	if pcfg.Cfg.Images.DeleteImageFiles {
		err = files.DeleteImage(filename)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(types.APIResponse{
				Body:  nil,
				Error: err.Error(),
			})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Deleted")
}
