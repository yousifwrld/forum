package main

import (
	"fmt"
	"log"
	"net/http"

	forum "forum/funcs"
)

func init() {
	err := forum.InitDatabase()
	if err != nil {
		log.Println(err)
		return
	}
}

func main() {
	http.HandleFunc("/", forum.RootHandler)
	http.HandleFunc("/register", forum.RegisterHandler)
	http.HandleFunc("/login", forum.LoginHandler)
	http.HandleFunc("/home", forum.HomeHandler)
	fmt.Println("server is listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
