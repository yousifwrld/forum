package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	forum "forum/funcs"
)

func init() {
	fmt.Println("initializing...")
	start := time.Now()
	err := forum.InitDatabase()
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

	http.Handle("/", forum.AuthMiddleware(http.HandlerFunc(forum.RootHandler)))
	http.HandleFunc("/register", forum.RegisterHandler)
	http.HandleFunc("/login", forum.LoginHandler)
	http.HandleFunc("/logout", forum.LogoutHandler)
	http.HandleFunc("/home", forum.HomeHandler)
	http.HandleFunc("/create-post", forum.CreatePostHandler)
	http.HandleFunc("/home/post/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/comment") {
			forum.CommentHandler(w, r)
		} else {
			forum.PostDetailHandler(w, r)
		}
	})

	fmt.Println("server is listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
