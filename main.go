package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"forum/db"
	forum "forum/funcs"
	"forum/oAuth"
)

func init() {
	fmt.Println("initializing...")
	start := time.Now()
	err := db.InitDatabase()
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("initializing took:")
	fmt.Println(time.Since(start))
}

func main() {
	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("./static")) // Ensure path is correct
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", forum.HomeHandler)
	http.HandleFunc("/register", forum.RegisterHandler)
	http.HandleFunc("/login", forum.LoginHandler)
	http.HandleFunc("/login/github", oAuth.GithubloginHandler)
	http.HandleFunc("/githubcallback", oAuth.GithubCallbackHandler)
	http.HandleFunc("/login/google", oAuth.GoogleLoginHandler)
	http.HandleFunc("/googlecallback", oAuth.GoogleCallbackHandler)
	http.HandleFunc("/logout", forum.LogoutHandler)
	http.HandleFunc("/help", forum.Help)
	http.Handle("/user-info", forum.AuthMiddleware(http.HandlerFunc(forum.UserInfo)))
	http.Handle("/reaction", forum.AuthMiddleware(http.HandlerFunc(forum.ReactionHandler)))
	http.HandleFunc("/filter", forum.FilterHandler)
	http.Handle("/create-post", forum.AuthMiddleware(http.HandlerFunc(forum.CreatePostHandler)))
	http.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/comment") {
			forum.AuthMiddleware(http.HandlerFunc(forum.CommentHandler)).ServeHTTP(w, r)
		} else {
			forum.PostDetailHandler(w, r)
		}
	})

	fmt.Println("server is listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)

	defer db.Database.Close()
}
