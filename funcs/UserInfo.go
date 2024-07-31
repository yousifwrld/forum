package forum

import (
	"net/http"
)

func UserInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorPages(w, r, "405", http.StatusMethodNotAllowed)
		return
	}
	RenderTemplate(w, "templates/userinfo.html", nil)

}
