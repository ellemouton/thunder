package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseGlob("assets/*.html"))

func main() {
	s, err := newState()
	if err != nil {
		log.Fatalf("newState: %s", err)
	}

	r := newRouter(s)

	fmt.Println("Serving on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
