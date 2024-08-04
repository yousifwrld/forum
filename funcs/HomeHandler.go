package forum

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"sort"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	type PageData struct {
		Posts              []Post
		IsLoggedIn         bool
		FilteredCategories []Category
	}
	if r.URL.Path != "/" {
		ErrorPages(w, r, "404", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
		return
	}

	isLoggedIn := false

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

	if userID != 0 {
		isLoggedIn = true
	}

	posts, err := getPosts()
	if err != nil {
		log.Println("Error getting posts:", err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	//get the categories for the filters
	var filteredCategories []Category
	rows, err := database.Query(`SELECT categoryID, name FROM category`)
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
	tmpl := template.Must(template.New("home.html").Funcs(template.FuncMap{
		"joinAndTrim": joinAndTrim,
	}).ParseFiles("templates/home.html"))

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

}
func getPosts() ([]Post, error) {
	rows, err := database.Query(`
        SELECT p.postID, p.title, p.content, p.created_at, u.username 
        FROM post p
        JOIN user u ON p.userID = u.userID
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.Username); err != nil {
			return nil, err
		}
		post.FormattedCreatedAt = post.CreatedAt.Format("2006-01-02 15:04")

		// Fetch categories for this post
		categories, err := GetCategoriesForPost(post.ID)
		if err != nil {
			return nil, err
		}
		// Get likes and dislikes for this post
		likes, dislikes, err := getLikesDislikes(post.ID)
		if err != nil {
			return nil, err
		}
		//save likes, dislikes and categories for each post
		post.Likes = likes
		post.Dislikes = dislikes
		post.Categories = categories

		// Append post to posts
		posts = append(posts, post)
	}
	// sort from latest to oldest based on created_at
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.After(posts[j].CreatedAt)
	})
	return posts, nil
}

func getLikesDislikes(postID int) (int, int, error) {
	var likes, dislikes int
	err := database.QueryRow(`
		SELECT likes, dislikes
		FROM post
		WHERE postID = ?
	`, postID).Scan(&likes, &dislikes)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, nil
		}
		return 0, 0, err
	}
	return likes, dislikes, nil
}
