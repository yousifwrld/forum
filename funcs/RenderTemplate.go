package forum

import (
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

	// Set content type to HTML
	w.Header().Set("Content-Type", "text/html")

	// Execute the template with the provided data
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}
