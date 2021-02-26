package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"github.com/ellemouton/thunder/blogs"
	blogs_db "github.com/ellemouton/thunder/blogs/db"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

// TODO(elle): Do better protection
const pssedstr = "$2a$10$LJLTJe93TOhoRhD.ZTyi6.Crdskdx4XJdbf1IueI/7BN9wyJNc6BG"

func newRouter(s *State) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", s.homeHandler).Methods("GET")
	r.HandleFunc("/blog", s.blogHandler).Methods("GET")
	r.HandleFunc("/blog/view/{id}", s.blogShowHandler).Methods("GET")
	r.HandleFunc("/gallery", s.galleryHandler).Methods("GET")
	r.HandleFunc("/projects", s.projectsHandler).Methods("GET")
	r.PathPrefix("/images").Handler(http.StripPrefix("/images", http.FileServer(http.Dir("assets/images"))))
	r.PathPrefix("/css").Handler(http.StripPrefix("/css", http.FileServer(http.Dir("assets/css"))))
	r.HandleFunc("/cv", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "assets/other/CV_Elle_Mouton.pdf") }).Methods("GET")

	// Hidden endpoints
	r.HandleFunc("/blog/add", s.newBlogTemplateHandler).Methods("GET")
	r.HandleFunc("/blog/add", s.newBlogHandler).Methods("POST")
	r.HandleFunc("/blog/edit/{id}", s.editBlogHandler).Methods("GET")
	r.HandleFunc("/blog/edit/{id}", s.saveEditBlogHandler).Methods("POST")

	return r
}

func (s *State) homeHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *State) saveEditBlogHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)
	r.ParseForm()

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	title := r.FormValue("title")
	abstract := r.FormValue("abstract")
	content := r.FormValue("content")
	password := r.FormValue("password")

	err = bcrypt.CompareHashAndPassword([]byte(pssedstr), []byte(password))
	if err != nil {
		http.Error(w, fmt.Sprintf("incorrect password, %s", err), http.StatusInternalServerError)
		return
	}

	err = blogs_db.UpdateBlog(ctx, s.GetDB(), id, title, abstract, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/blog/view/%d", id), http.StatusSeeOther)
}

func (s *State) editBlogHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	info, err := blogs_db.LookupInfo(ctx, s.GetDB(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content, err := blogs_db.LookupContent(ctx, s.GetDB(), info.ContentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := struct {
		ID       int64
		Name     string
		Abstract string
		Date     time.Time
		Content  string
	}{
		ID:       id,
		Name:     info.Name,
		Abstract: info.Description,
		Content:  content.Text,
	}

	err = templates.ExecuteTemplate(w, "edit_blog.html", c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *State) blogShowHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	info, err := blogs_db.LookupInfo(ctx, s.GetDB(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content, err := blogs_db.LookupContent(ctx, s.GetDB(), info.ContentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
	var buf bytes.Buffer
	if err := md.Convert([]byte(content.Text), &buf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := struct {
		Name     string
		Abstract string
		Date     time.Time
		Content  template.HTML
	}{
		Name:     info.Name,
		Abstract: info.Description,
		Date:     info.CreatedAt,
		Content:  template.HTML(string(buf.Bytes())),
	}

	err = templates.ExecuteTemplate(w, "blog.html", c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *State) blogHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	bls, err := blogs_db.ListAllInfoRev(ctx, s.GetDB())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := struct {
		Payload []*blogs.Info
	}{
		Payload: bls,
	}

	err = templates.ExecuteTemplate(w, "blogs.html", c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *State) newBlogTemplateHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "add_blog.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *State) newBlogHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	r.ParseForm()

	title := r.FormValue("title")
	abstract := r.FormValue("abstract")
	content := r.FormValue("content")
	password := r.FormValue("password")

	err := bcrypt.CompareHashAndPassword([]byte(pssedstr), []byte(password))
	if err != nil {
		http.Error(w, fmt.Sprintf("incorrect password, %s", err), http.StatusInternalServerError)
		return
	}

	_, err = blogs_db.Create(ctx, s.GetDB(), title, abstract, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.homeHandler(w, r)
}

func (s *State) galleryHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "gallery.html", nil)
	if err != nil {
		log.Println(err)
	}
}

func (s *State) projectsHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "projects.html", nil)
	if err != nil {
		log.Println(err)
	}
}
