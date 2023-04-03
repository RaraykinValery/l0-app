package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/RaraykinValery/l0/cache"
	"github.com/RaraykinValery/l0/models"
	"github.com/RaraykinValery/l0/subscriber"
)

var templates_paths = []string{
	filepath.Join("templates", "index.html"),
	filepath.Join("templates", "order.html"),
}

var templates = template.Must(template.ParseFiles(templates_paths...))

func orderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := templates.ExecuteTemplate(w, "index", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == http.MethodPost {
		var order models.Order

		uuid := r.FormValue("uuid")

		order, ok := cache.GetOrderFromCache(uuid)
		if ok == false {
			log.Printf("The order is not in cache")
		}

		templates.ExecuteTemplate(w, "index", order)
	}

}

func main() {
	go subscriber.StartSubscriber()

	fileServer := http.FileServer(http.Dir("./static/"))

	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	http.HandleFunc("/", orderHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
