package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gomarkdown/markdown"
	"github.com/gorilla/mux"
)

var templates = template.Must(template.ParseGlob("assets/*.html"))

func main() {
	r := newRouter()

	fmt.Println("Serving on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func newRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/blog", blogHandler).Methods("GET")
	r.HandleFunc("/gallery", galleryHandler).Methods("GET")
	r.HandleFunc("/projects", projectsHandler).Methods("GET")
	r.PathPrefix("/images").Handler(http.StripPrefix("/images", http.FileServer(http.Dir("assets/images"))))
	r.HandleFunc("/cv", cvHandler).Methods("GET")

	return r
}

func cvHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/other/CV_Elle_Mouton.pdf")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		log.Println(err)
	}
}

type Content struct {
	Text template.HTML
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	md := []byte("## markdown document\n" +
		"![testing](/images/well.png)\n" +
		"blah `blah` blah")

	output := markdown.ToHTML(md, nil, nil)

	c := Content{
		Text: template.HTML(output),
	}

	err := templates.ExecuteTemplate(w, "blog.html", c)
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
