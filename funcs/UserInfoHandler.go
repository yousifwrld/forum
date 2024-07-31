package forum

import (
	"html/template"
	"log"
	"net/http"
	"slices"
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
	posts, err := getPostsByUserID(userID)
	if err != nil {
		log.Println(err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	//get the liked posts for the user
	likedPosts, err := getLikedPostsByUserID(userID)
	if err != nil {
		log.Println(err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Username:   username,
		Posts:      posts,
		LikedPosts: likedPosts,
	}

	//to be able to use the joinAndTrim function in the html
	tmpl := template.Must(template.New("userinfo.html").Funcs(template.FuncMap{
		"joinAndTrim": joinAndTrim,
	}).ParseFiles("templates/userinfo.html"))

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

}

func getPostsByUserID(userID int) ([]Post, error) {

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

func getLikedPostsByUserID(userID int) ([]Post, error) {

	//start a transaction
	tx, err := database.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	rows, err := tx.Query(`SELECT postID FROM reaction WHERE userID = ? AND is_like = true`, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	// get the posts
	var posts []Post

	for rows.Next() {
		var postID int
		if err := rows.Scan(&postID); err != nil {
			log.Println(err)
			return nil, err
		}

		// Now we query the details for the post
		postRow := tx.QueryRow(`
		SELECT p.postID, p.title, p.content, p.likes, p.dislikes, p.created_at, u.username
		FROM post p
		JOIN user u ON p.userID = u.userID
		WHERE p.postID = ?
		`, postID)

		var post Post
		if err := postRow.Scan(&post.ID, &post.Title, &post.Content, &post.Likes, &post.Dislikes, &post.CreatedAt, &post.Username); err != nil {
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

	// Check if there were any errors during iteration
	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	// in case of no posts liked by the user, set the posts to an empty slice
	if posts == nil {
		posts = []Post{}
	}

	//reverse the order of the posts to display them based on the latest liked posts
	slices.Reverse(posts)

	return posts, nil
}
