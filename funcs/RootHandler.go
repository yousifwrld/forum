package forum

import "net/http"

func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorPages(w, r, "400", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, "/register", http.StatusFound)
}
