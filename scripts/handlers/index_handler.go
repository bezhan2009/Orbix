package handlers

import (
	"goCmd/cmd/commands"
	"html/template"
	"log"
	"net/http"
)

// IndexHandler Handler for rendering index.html
func IndexHandler(w http.ResponseWriter,
	r *http.Request) {
	setHeaders(w)
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Template error: %v", err)
		log.Printf("Changing directory for finding templates")

		err = commands.ChangeDirectory("..")
		if err != nil {
			log.Printf("Change directory error: %v", err)
		}
		tmpl, err = template.ParseFiles("templates/index.html")
		if err != nil {
			log.Printf("Template error: %v", err)
			http.Error(w, "Unable to load template", http.StatusInternalServerError)
			return
		}
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return
	}
}
