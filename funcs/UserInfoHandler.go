package forum

import (
	"database/sql"
	"log"
	"net/http"
	"sort"
)

func UserInfo(w http.ResponseWriter, r *http.Request) {

	type PageData struct {
		Username   string
		Posts      []Post
		LikedPosts []Post
	}
	if r.Method != http.MethodGet {
		ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
		return
	}

	//get the user id from request context
	userID := r.Context().Value(userIDKey).(int)

	//get the username from the database
	var username string
	err := database.QueryRow(`SELECT username FROM user WHERE userID = ?`, userID).Scan(&username)
	if err != nil {
		log.Println(err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	//get the posts created by the user
	posts, err := getPostsByID(userID)
	if err != nil {
		log.Println(err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	//get the liked posts for the user
	var likedPosts []Post

	//get the user's liked posts from the database
	rows, err := database.Query(`SELECT postID FROM reaction WHERE userID = ? AND is_like = true`, userID)
	if err != nil {
		//if no results this means that the user hasnt liked any post before
		if err == sql.ErrNoRows {
			//set the liked posts to an empty slice
			likedPosts = []Post{}
		} else {
			log.Println(err)
			ErrorPages(w, r, "500", http.StatusInternalServerError)
			return
		}
	}

	for rows.Next() {

		var postID int
		err := rows.Scan(&postID)
		if err != nil {
			log.Println(err)
			ErrorPages(w, r, "500", http.StatusInternalServerError)
			return
		}

		//get the post from the database
		posts, err := getPostsByID(postID)
		if err != nil {
			log.Println(err)
			ErrorPages(w, r, "500", http.StatusInternalServerError)
			return
		}
		likedPosts = posts
	}

	data := PageData{
		Username:   username,
		Posts:      posts,
		LikedPosts: likedPosts,
	}

	RenderTemplate(w, "templates/userinfo.html", data)

}

func getPostsByID(userID int) ([]Post, error) {

	//start a transaction
	tx, err := database.Begin()
	if err != nil {
		return nil, err
	}

	//
	rows, err := tx.Query(`SELECT postID, title, content, likes, dislikes, created_at FROM post WHERE userID = ?`, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// a defered function that will rollback the transaction in case of an error, or commit it otherwise
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	//get the user's posts from the database
	var posts []Post

	for rows.Next() {
		var post Post
		//scan the rows into the post struct
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Likes, &post.Dislikes, &post.CreatedAt); err != nil {
			log.Println(err)
			return nil, err
		}

		//get the categories for each post
		categories, err := GetCategoriesForPost(post.ID)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		//save the categories for each post, and format the time
		post.Categories = categories
		post.FormattedCreatedAt = post.CreatedAt.Format("2006-01-02 15:04")

		posts = append(posts, post)
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// in case of no posts created by the user, set the posts to an empty slice
	if posts == nil {
		posts = []Post{}
	}

	//sort from latest to oldest based on created_at
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.After(posts[j].CreatedAt)
	})
	return posts, nil
}
