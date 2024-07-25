package forum

import (
	"fmt"
	"html/template"
	"net/http"
)

// RenderTemplate renders the HTML template with the given data
func RenderTemplate(w http.ResponseWriter, templateFile string, data interface{}) {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	// Execute the template with the provided data
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println(templateFile)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}
