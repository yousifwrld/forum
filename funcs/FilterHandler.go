package forum

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func FilterHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		ErrorPages(w, r, "400", http.StatusBadRequest)
		return
	}

	categories := r.Form["filter"]

	var catID []int

	for _, id := range categories {
		intID, err := strconv.Atoi(id)
		if err != nil {
			log.Println(err)
			ErrorPages(w, r, "400", http.StatusBadRequest)
			return
		}
		catID = append(catID, intID)
	}

	statement := fmt.Sprintf(`SELECT postID from post_category WHERE categoryID = %d`, catID[0])
	for i := 1; i < len(catID); i++ {
		statement += fmt.Sprintf(`OR categoryID = %d`, catID[i])
	}
	rows, err := database.Query(statement)

	for rows.Next() {
		var postID int
		rows.Scan(&postID)
		database.Query(`SELECT title, content, likes, dislikes, created_at, userID WHERE postID = ?`, postID)
	}
}
