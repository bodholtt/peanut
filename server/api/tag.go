package api

import (
	"encoding/json"
	"log"
	"net/http"
	"peanutserver/database"
	"peanutserver/pcfg"
	"peanutserver/types"
	"strconv"
)

// API Route: /tag

func HandleTag(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	switch r.Method {
	case "GET":
		handleTagGET(w, r)
	case "POST":
		handleTagPOST(w, r)
	case "PUT":
		handleTagPUT(w, r)
	case "DELETE":
		handleTagDELETE(w, r)
	case "OPTIONS":
		handleTagOPTIONS(w, r)
	}
}

// handleTagGET - get a tag by its id
func handleTagGET(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: err.Error(),
		})
		return
	}

	tag, err := database.GetTag(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(types.APIResponse{
		Body:  tag,
		Error: "",
	})
}

func handleTagOPTIONS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", pcfg.Cfg.Client.Host)
}

func handleTagPOST(w http.ResponseWriter, r *http.Request) {

	tag, err := database.CreateTag("tagtest")
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
		Body:  tag,
		Error: "",
	})
}

func handleTagPUT(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Not implemented")
}

// handleTagDELETE - DELETE BY *** TAG ID ***
func handleTagDELETE(w http.ResponseWriter, r *http.Request) {

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

	err = database.DeleteTag(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(types.APIResponse{
			Body:  nil,
			Error: err.Error(),
		})
		return
	}

	log.Println("Deleted tag", id)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.APIResponse{
		Body:  "Deleted tag",
		Error: "",
	})
}
