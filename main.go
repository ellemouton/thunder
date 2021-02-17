package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var templates = template.Must(template.ParseGlob("assets/*.html"))

func main() {
	r := newRouter()

	fmt.Println("Serving on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/blog", blogHandler).Methods("GET")
	r.HandleFunc("/gallery", galleryHandler).Methods("GET")
	r.HandleFunc("/projects", projectsHandler).Methods("GET")
	return r
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		log.Println(err)
	}
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "blog.html", nil)
	if err != nil {
		log.Println(err)
	}
}

func galleryHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "gallery.html", nil)
	if err != nil {
		log.Println(err)
	}
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "projects.html", nil)
	if err != nil {
		log.Println(err)
	}
}
