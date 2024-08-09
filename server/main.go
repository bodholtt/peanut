package main

// MY MAIN.GO IS TO BLOW UP
// THEN ACT LIKE I DON'T KNOW NOBODY

import (
	"log"
	"net/http"
	"peanutserver/api"
	"peanutserver/auth"
	"peanutserver/database"
	"peanutserver/pcfg"
	"strconv"
)

func main() {

	err := pcfg.InitConfig("config.yml")
	if err != nil {
		panic(err)
	}

	database.Initialize()

	mux := http.NewServeMux()

	// file server (route: /static/)
	fs := http.FileServer(http.Dir(pcfg.Cfg.Images.RootLocation))
	log.Println("Serving files from:", http.Dir(pcfg.Cfg.Images.RootLocation))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// route: /posts
	mux.Handle("/posts", auth.RankMiddleware(http.HandlerFunc(api.HandlePosts), pcfg.Cfg.Permissions.ViewPosts))

	// route: /postsCount
	mux.HandleFunc("/postCount", api.HandlePostCount)

	// route: /post
	mux.HandleFunc("OPTIONS /post", api.HandlePostOPTIONS)
	mux.HandleFunc("OPTIONS /post/{id}", api.HandlePostOPTIONS)
	mux.Handle("POST /post", auth.RankMiddleware(http.HandlerFunc(api.HandlePostPOST), pcfg.Cfg.Permissions.CreatePosts))
	mux.Handle("GET /post/{id}", auth.RankMiddleware(http.HandlerFunc(api.HandlePostGET), pcfg.Cfg.Permissions.ViewPosts))
	// TODO: Allow users to edit their own posts, but restrict editing posts that are not their own.
	mux.Handle("PUT /post/{id}", auth.RankMiddleware(http.HandlerFunc(api.HandlePostPUT), pcfg.Cfg.Permissions.EditOthersPosts))
	// TODO: Restrict deletion of own posts and others' posts separately.
	mux.Handle("DELETE /post/{id}", auth.RankMiddleware(http.HandlerFunc(api.HandlePostDELETE), pcfg.Cfg.Permissions.DeleteOwnPosts))

	// route: /login
	mux.HandleFunc("POST /login", api.HandleLogin)
	mux.HandleFunc("OPTIONS /login", api.HandleAccountsOPTIONS)

	// route: /signup
	mux.Handle("POST /signup", auth.RankMiddleware(http.HandlerFunc(api.HandleSignup), pcfg.Cfg.Permissions.SignUp))
	mux.HandleFunc("OPTIONS /signup", api.HandleAccountsOPTIONS)

	// route: /createUser
	mux.Handle("POST /createUser", auth.RankMiddleware(http.HandlerFunc(api.HandleCreateUser), pcfg.Cfg.Permissions.CreateUsers))
	mux.HandleFunc("OPTIONS /createUser", api.HandleAccountsOPTIONS)

	// route: /user
	mux.HandleFunc("/user", api.HandleUser)
	mux.HandleFunc("/user/{id}", api.HandleUser)

	// route: /tag
	mux.HandleFunc("/tag", api.HandleTag)
	mux.HandleFunc("/tag/{id}", api.HandleTag)

	mux.HandleFunc("/post/{id}/tags", api.HandlePostTags)

	log.Println("listening on port", pcfg.Cfg.Server.Port)

	err = http.ListenAndServe(":"+strconv.Itoa(pcfg.Cfg.Server.Port), mux)
	if err != nil {
		panic(err)
	}
}
