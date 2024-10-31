package forum

import (
	"encoding/base64"
	"forum/db"
	"html/template"
	"log"
	"net/http"
	"sort"
)

type PostWithComments struct {
	Post     Post
	Comments []Comment
}

func UserInfo(w http.ResponseWriter, r *http.Request) {

	type PageData struct {
		Username      string
		Posts         []Post
		LikedPosts    []Post
		DislikedPosts []Post
		Comments      map[int]PostWithComments
	}
	if r.Method != http.MethodGet {
		ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
		return
	}

	//get the user id from request context
	userID := r.Context().Value(userIDKey).(int)

	//get the username from the database
	var username string
	err := db.Database.QueryRow(`SELECT username FROM user WHERE userID = ?`, userID).Scan(&username)
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
	likedPosts, dislikedPosts, err := GetLikedAndDislikedPosts(userID)
	if err != nil {
		log.Println(err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	//get the comments created by the user
	postAndComments, err := getUserCommentsOnPost(userID)
	if err != nil {
		log.Println(err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Username:      username,
		Posts:         posts,
		LikedPosts:    likedPosts,
		DislikedPosts: dislikedPosts,
		Comments:      postAndComments,
	}

	//to be able to use the joinAndTrim function in the html
	tmpl, err := template.New("userinfo.html").Funcs(template.FuncMap{
		"joinAndTrim": joinAndTrim,
	}).ParseFiles("templates/userinfo.html")

	if err != nil {
		log.Println(err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

}

func getPostsByUserID(userID int) ([]Post, error) {

	//start a transaction
	tx, err := db.Database.Begin()
	if err != nil {
		return nil, err
	}

	//
	rows, err := tx.Query(`SELECT postID, title, content, image, likes, dislikes, created_at FROM post WHERE userID = ?`, userID)
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
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Image, &post.Likes, &post.Dislikes, &post.CreatedAt); err != nil {
			log.Println(err)
			return nil, err
		}

		//get the categories for each post
		categories, err := GetCategoriesForPost(post.ID)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		//base64 encode the image to display it in the HTML
		post.Base64Image = base64.StdEncoding.EncodeToString(post.Image)
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

func GetLikedAndDislikedPosts(userID int) ([]Post, []Post, error) {
	// Start a transaction
	tx, err := db.Database.Begin()
	if err != nil {
		return nil, nil, err
	}

	// Rollback the transaction on error or commit on success
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Query to get post IDs along with their like or dislike status
	rows, err := tx.Query(`
		SELECT postID, is_like, is_dislike 
		FROM reaction 
		WHERE userID = ? AND (is_like = true OR is_dislike = true)
	`, userID)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	defer rows.Close()

	// Prepare slices for liked and disliked posts
	var likedPosts []Post
	var dislikedPosts []Post

	for rows.Next() {
		var postID int
		var isLike, isDislike bool
		if err := rows.Scan(&postID, &isLike, &isDislike); err != nil {
			log.Println(err)
			return nil, nil, err
		}

		// Query the details for each post
		postRow := tx.QueryRow(`
			SELECT p.postID, p.title, p.content, p.image, p.likes, p.dislikes, p.created_at, u.username
			FROM post p
			JOIN user u ON p.userID = u.userID
			WHERE p.postID = ? ORDER BY p.created_at DESC
		`, postID)

		// Prepare a Post struct to hold the post details
		var post Post
		if err := postRow.Scan(&post.ID, &post.Title, &post.Content, &post.Image, &post.Likes, &post.Dislikes, &post.CreatedAt, &post.Username); err != nil {
			log.Println(err)
			return nil, nil, err
		}

		// Get the categories for each post
		categories, err := GetCategoriesForPost(post.ID)
		if err != nil {
			log.Println(err)
			return nil, nil, err
		}

		// Base64 encode the image and set additional fields
		post.Base64Image = base64.StdEncoding.EncodeToString(post.Image)
		post.Categories = categories
		post.FormattedCreatedAt = post.CreatedAt.Format("2006-01-02 15:04")

		// Append post to likedPosts or dislikedPosts based on its status
		if isLike {
			likedPosts = append(likedPosts, post)
		} else if isDislike {
			dislikedPosts = append(dislikedPosts, post)
		}
	}

	// Check for any errors during iteration
	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, nil, err
	}

	return likedPosts, dislikedPosts, nil
}

func getUserCommentsOnPost(userID int) (map[int]PostWithComments, error) {
	// Start a transaction
	tx, err := db.Database.Begin()
	if err != nil {
		return nil, err
	}

	// Prepare the map to group comments by post
	postAndComments := make(map[int]PostWithComments)

	// Query to get user's comments
	rows, err := tx.Query(`SELECT commentID, comment, likes, dislikes, created_at, postID FROM comment WHERE userID = ? ORDER BY created_at DESC`, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Defer a function to commit or rollback the transaction
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Process each comment
	for rows.Next() {
		var comment Comment
		var postID int
		// Scan comment data including postID
		if err := rows.Scan(&comment.CommentID, &comment.Content, &comment.Likes, &comment.Dislikes, &comment.CreatedAt, &postID); err != nil {
			log.Println(err)
			return nil, err
		}

		// Format the comment's creation date
		comment.FormattedCreatedAt = comment.CreatedAt.Format("2006-01-02 15:04")

		// Retrieve the post and comments slice from the map (or initialize if not present)
		postWithComments, exists := postAndComments[postID]
		if !exists {
			// Fetch the post data if this is the first comment for this postID
			post, err := getPostByPostID(postID)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			postWithComments = PostWithComments{
				Post:     post,
				Comments: []Comment{},
			}
		}

		// Append the new comment to the Comments slice
		postWithComments.Comments = append(postWithComments.Comments, comment)

		// Reassign the modified struct back to the map
		postAndComments[postID] = postWithComments
	}

	// Return the map of posts and comments
	return postAndComments, nil
}

func getPostByPostID(postID int) (Post, error) {

	//start a transaction
	tx, err := db.Database.Begin()
	if err != nil {
		return Post{}, err
	}

	//get the post from the database
	var post Post
	err = tx.QueryRow(`SELECT postID, title, content, image, likes, dislikes, created_at FROM post WHERE postID = ?`, postID).Scan(&post.ID, &post.Title, &post.Content, &post.Image, &post.Likes, &post.Dislikes, &post.CreatedAt)
	if err != nil {
		log.Println(err)
		return Post{}, err
	}

	//get the categories for the post
	categories, err := GetCategoriesForPost(post.ID)
	if err != nil {
		log.Println(err)
		return Post{}, err
	}

	//base64 encode the image to display it in the HTML
	post.Base64Image = base64.StdEncoding.EncodeToString(post.Image)
	//save the categories for the post, and format the time
	post.Categories = categories
	post.FormattedCreatedAt = post.CreatedAt.Format("2006-01-02 15:04")

	return post, nil
}
