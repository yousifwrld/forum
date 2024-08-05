package forum

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
)

func FilterHandler(w http.ResponseWriter, r *http.Request) {

	type PageData struct {
		Posts            []Post
		IsLoggedIn       bool
		FilterCategories []Category
	}

	if r.Method != http.MethodPost {
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

	err = r.ParseForm()
	if err != nil {
		log.Println(err)
		ErrorPages(w, r, "400", http.StatusBadRequest)
		return
	}

	categories := r.Form["filter"]

	if len(categories) == 0 {
		ErrorPages(w, r, "400", http.StatusBadRequest)
		return
	}

	// Convert category IDs to integers
	var catID []int
	for _, id := range categories {
		intID, err := strconv.Atoi(id)
		if err != nil {
			log.Println(err)
			ErrorPages(w, r, "400", http.StatusBadRequest)
		}
		catID = append(catID, intID)
	}

	// Create SQL statement for filtering posts by category
	statement := fmt.Sprintf("SELECT DISTINCT postID FROM post_categories WHERE categoryID = %d", catID[0])

	for i := 1; i < len(catID); i++ {
		statement += fmt.Sprintf(" OR categoryID = %d", catID[i])
	}

	rows, err := database.Query(statement)
	if err != nil {
		log.Println(err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var postID int
		if err := rows.Scan(&postID); err != nil {
			log.Println(err)
			ErrorPages(w, r, "500", http.StatusInternalServerError)
		}

		var post Post
		err := database.QueryRow(
			`SELECT p.title, p.content, p.created_at, p.likes, p.dislikes, u.username
			FROM post p
			JOIN user u ON p.userID = u.userID
			WHERE p.postID = ?`, postID).Scan(&post.Title, &post.Content, &post.CreatedAt, &post.Likes, &post.Dislikes, &post.Username)
		if err != nil {
			log.Println(err)
			ErrorPages(w, r, "500", http.StatusInternalServerError)
		}

		post.ID = postID
		post.FormattedCreatedAt = post.CreatedAt.Format("2006-01-02 15:04")

		categories, err := GetCategoriesForPost(postID)
		if err != nil {
			log.Println(err)
			ErrorPages(w, r, "500", http.StatusInternalServerError)

		}
		post.Categories = categories

		posts = append(posts, post)
	}

	// Sort filtered posts by creation date
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.After(posts[j].CreatedAt)
	})

	//get the categories for the filters
	var filterCategories []Category
	rows, err = database.Query(`SELECT categoryID, name FROM category`)
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
		filterCategories = append(filterCategories, category)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error after iterating rows:", err)
		ErrorPages(w, r, "500", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Posts:      posts,
		IsLoggedIn: isLoggedIn,

		FilterCategories: filterCategories,
	}
	//to be able to use the joinAndTrim function in the html
	tmpl, err := template.New("filter.html").Funcs(template.FuncMap{
		"joinAndTrim": joinAndTrim,
	}).ParseFiles("templates/filter.html")

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
