package forum

import (
	"log"
	"net/http"
	"text/template"
)

// renderTemplate renders the given template with the provided data
func RenderTemplate(w http.ResponseWriter, templatePath string, data any) {
	temp, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Println(err)
		http.Error(w, "500: Internal Server Error!", http.StatusInternalServerError)
		return
	}
	err = temp.Execute(w, data)
	if err != nil {
		log.Println(err)
		http.Error(w, "500: Internal Server Error!", http.StatusInternalServerError)
		return
	}
}
