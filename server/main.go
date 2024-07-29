package main

// MY MAIN.GO IS TO BLOW UP
// THEN ACT LIKE I DON'T KNOW NOBODY

import (
	"net/http"
	"peanutserver/api"
	"peanutserver/database"
)

func main() {
	database.Initialize()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/post", api.HandlePost)
	http.HandleFunc("/post/", api.HandlePost)
	http.HandleFunc("/posts", api.HandlePosts)
	http.HandleFunc("/postCount", api.HandlePostCount)

	http.HandleFunc("/user", api.HandleUser)
	http.HandleFunc("/user/", api.HandleUser)

	http.ListenAndServe(":8080", nil)
}
