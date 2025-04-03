package handlers

import (
	"fmt"
	"goCmd/cmd/commands"
	"html/template"
	"log"
	"net/http"
	"os"
)

var templateDir string
var userDir, _ = os.Getwd()

// IndexHandler Handler for rendering index.html
func IndexHandler(w http.ResponseWriter,
	r *http.Request) {
	setHeaders(w)

	tmpl, err := template.ParseFiles("templates\\index.html")
	if err != nil {
		log.Printf("Template error: %v\n\tDir: %s", err, userDir)
		log.Printf("Changing directory for finding templates")

		err = commands.ChangeDirectory("..")
		if err != nil {
			log.Printf("Change directory error: %v", err)
			return
		}
		tmpl, err = template.ParseFiles("templates\\index.html")
		if err != nil {
			log.Printf("Template error: %v\n\tDir: %s", err, userDir)

			err = commands.ChangeDirectory(userDir)
			if err != nil {
				log.Printf("Change directory error: %v", err)
				return
			}

			tmplFile := fmt.Sprintf("%s\\templates\\index.html", userDir)

			tmpl, err = template.ParseFiles(tmplFile)
			if err == nil {
				err = tmpl.Execute(w, nil)
				if err != nil {
					log.Printf("Template error: %v\n\tDir: %s\n\tTmpl dir: %s", err, userDir, tmplFile)
					return
				}

				return
			}

			log.Printf("Template error: %v", err)
			http.Error(w, "Unable to load template", http.StatusInternalServerError)
			return
		}

		//userDir, _ = os.Getwd()
	}

	templateDir = userDir + "templates"

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Template error: %v", err)
		return
	}
}
