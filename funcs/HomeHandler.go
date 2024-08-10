package forum

import (
	"database/sql"
	"encoding/base64"
	"forum/db"
	"html/template"
	"log"
	"net/http"
	"sort"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	//struct with all the data we need for the page
	type PageData struct {
		Posts              []Post
		IsLoggedIn         bool
		FilteredCategories []Category
	}

	//validate path and method type
	if r.URL.Path != "/" {
		ErrorPages(w, r, "404", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
		return
	}

	//check if the user is logged in or no to display correct buttons (login/logout)
	isLoggedIn := false

	//get the sessionID from the cookie if available
	sessionID, err := GetSessionFromCookie(r)
	if err != nil {
		if err == http.ErrNoCookie {
			isLoggedIn = false
		} else {
			log.Println("Error getting session cookie:", err)
			ErrorPages(w, r, "500", http.StatusInternalServerError)
			return
		}
	}

	//if a session was available get the userID associated with the session
	userID, err := GetIDFromSession(sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			isLoggedIn = false
		} else {
			log.Println("Error getting userID from session:", err)
			ErrorPages(w, r, "500", http.StatusInternalServerError)
			return
		}
	}

	//if a userID was retrieved then the user is logged in
	if userID != 0 {
		isLoggedIn = true
	}

	//get all the posts for the home page
	posts, err := getPosts()
	if err != nil {
		log.Println("Error getting posts:", err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	//get the categories for the filters option

	//select all of the categories from the database
	var filteredCategories []Category
	rows, err := db.Database.Query(`SELECT categoryID, name FROM category`)
	if err != nil {
		log.Println("Error querying categories:", err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.CategoryID, &category.Name)
		if err != nil {
			log.Println("Error scanning category:", err)
			ErrorPages(w, r, "500", http.StatusInternalServerError)
			return
		}
		//append the category to the array, to loop over in html
		filteredCategories = append(filteredCategories, category)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error after iterating rows:", err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Posts:              posts,
		IsLoggedIn:         isLoggedIn,
		FilteredCategories: filteredCategories,
	}

	//to be able to use the joinAndTrim function in the html
	tmpl, err := template.New("home.html").Funcs(template.FuncMap{
		"joinAndTrim": joinAndTrim,
	}).ParseFiles("templates/home.html")

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
func getPosts() ([]Post, error) {
	//select all the posts with all the info from the database
	rows, err := db.Database.Query(`
        SELECT p.postID, p.title, p.content, p.image, p.created_at, u.username 
        FROM post p
        JOIN user u ON p.userID = u.userID
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	//loop over the result rows and append each post to the posts slice
	for rows.Next() {
		var post Post
		// Scan the result rows into the post struct
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Image, &post.CreatedAt, &post.Username); err != nil {
			return nil, err
		}

		// Base64 encode the image to display it in the HTML
		post.Base64Image = base64.StdEncoding.EncodeToString(post.Image)
		//format the time into a more readable format
		post.FormattedCreatedAt = post.CreatedAt.Format("2006-01-02 15:04")

		// Fetch categories associated with this post
		categories, err := GetCategoriesForPost(post.ID)
		if err != nil {
			return nil, err
		}
		// Get likes and dislikes and comments count for this post
		likes, dislikes, comments, err := getLikesDislikesComments(post.ID)
		if err != nil {
			return nil, err
		}
		//save likes, dislikes, comments and categories for each post
		post.Likes = likes
		post.Dislikes = dislikes
		post.Categories = categories
		post.Comments = comments

		// Append post to posts
		posts = append(posts, post)
	}
	// sort from latest to oldest based on created_at
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.After(posts[j].CreatedAt)
	})
	return posts, nil
}

func getLikesDislikesComments(postID int) (int, int, int, error) {
	var likes, dislikes int
	err := db.Database.QueryRow(`
		SELECT likes, dislikes
		FROM post
		WHERE postID = ?
	`, postID).Scan(&likes, &dislikes)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, 0, nil
		}
		return 0, 0, 0, err
	}

	var comments int
	err = db.Database.QueryRow(`SELECT COUNT(*) FROM comment WHERE postID = ?`, postID).Scan(&comments)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, 0, nil
		}
		return 0, 0, 0, err
	}

	return likes, dislikes, comments, nil
}
