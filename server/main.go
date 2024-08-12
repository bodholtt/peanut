package main

// MY MAIN.GO IS TO BLOW UP
// THEN ACT LIKE I DON'T KNOW NOBODY

import (
	"log"
	"net/http"
	"peanutserver/api"
	"peanutserver/api/user"
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
	mux.Handle("/posts", auth.RankMiddleware(http.HandlerFunc(api.HandlePosts), pcfg.Perms.ViewPosts))

	// route: /postsCount
	mux.HandleFunc("/postCount", api.HandlePostCount)

	// route: /post
	mux.HandleFunc("OPTIONS /post", api.HandlePostOPTIONS)
	mux.HandleFunc("OPTIONS /post/{id}", api.HandlePostOPTIONS)
	mux.Handle("POST /post", auth.RankMiddleware(http.HandlerFunc(api.HandlePostPOST), pcfg.Perms.CreatePosts))
	mux.Handle("GET /post/{id}", auth.RankMiddleware(http.HandlerFunc(api.HandlePostGET), pcfg.Perms.ViewPosts))
	// TODO: Allow users to edit their own posts, but restrict editing posts that are not their own.
	mux.Handle("PUT /post/{id}", auth.RankMiddleware(http.HandlerFunc(api.HandlePostPUT), pcfg.Perms.EditOthersPosts))
	// TODO: Restrict deletion of own posts and others' posts separately.
	mux.Handle("DELETE /post/{id}", auth.RankMiddleware(http.HandlerFunc(api.HandlePostDELETE), pcfg.Perms.DeleteOwnPosts))

	// route: /login
	mux.HandleFunc("POST /login", user.HandleLogin)
	mux.HandleFunc("OPTIONS /login", user.HandleAccountsOPTIONS)

	// route: /signup
	mux.Handle("POST /signup", auth.RankMiddleware(http.HandlerFunc(user.HandleSignup), pcfg.Perms.SignUp))
	mux.HandleFunc("OPTIONS /signup", user.HandleAccountsOPTIONS)

	// route: /user
	mux.HandleFunc("OPTIONS /user", user.HandleUserOPTIONS)
	mux.Handle("POST /user", auth.RankMiddleware(http.HandlerFunc(user.HandleUserPOST), pcfg.Perms.CreateUsers))
	mux.HandleFunc("GET /user/{id}", user.HandleUserGET)
	mux.HandleFunc("GET /user/{id}/permissions", user.CheckPermissions)

	// route: /tag
	mux.HandleFunc("/tag", api.HandleTag)
	mux.HandleFunc("/tag/{id}", api.HandleTag)

	log.Println("listening on port", pcfg.Cfg.Server.Port)

	err = http.ListenAndServe(":"+strconv.Itoa(pcfg.Cfg.Server.Port), mux)
	if err != nil {
		panic(err)
	}
}
