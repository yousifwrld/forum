package main

import (
	"fmt"
	"log"
	"net/http"
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

	http.HandleFunc("/", forum.RootHandler)
	http.HandleFunc("/register", forum.RegisterHandler)
	http.HandleFunc("/login", forum.LoginHandler)
	http.HandleFunc("/logout", forum.LogoutHandler)
	http.HandleFunc("/home", forum.HomeHandler)
	http.HandleFunc("/home/post-detail/", forum.PostDetailHandler)
	http.HandleFunc("/create-post", forum.CreatePostHandler)
	http.HandleFunc("/home/post-detail/comment", forum.CommentHandler)

	fmt.Println("server is listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
