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

	"github.com/ellemouton/thunder/blogs"
	blogs_db "github.com/ellemouton/thunder/blogs/db"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

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

/*
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

	//extensions := parser.CommonExtensions | parser.HardLineBreak
	extensions := parser.NoIntraEmphasis        // Ignore emphasis markers inside words
	extensions |= parser.Tables                 // Parse tables
	extensions |= parser.FencedCode             // Parse fenced code blocks
	extensions |= parser.Autolink               // Detect embedded URLs that are not explicitly marked
	extensions |= parser.Strikethrough          // Strikethrough text using ~~test~~
	extensions |= parser.SpaceHeadings          // Be strict about prefix heading rules
	extensions |= parser.HeadingIDs             // specify heading IDs  with {#id}
	extensions |= parser.BackslashLineBreak     // Translate trailing backslashes into line breaks
	extensions |= parser.DefinitionLists        // Parse definition lists
	extensions |= parser.LaxHTMLBlocks          // more in HTMLBlock, less in HTMLSpan
	extensions |= parser.NoEmptyLineBeforeBlock // no need for new line before a list
	extensions |= parser.HardLineBreak
	extensions |= parser.EmptyLinesBreakList
	parser := parser.NewWithExtensions(extensions)
	output := markdown.ToHTML([]byte(content.Text), parser, nil)

	c := struct {
		Name     string
		Abstract string
		Date     time.Time
		Content  template.HTML
	}{
		Name:     info.Name,
		Abstract: info.Description,
		Date:     info.CreatedAt,
		Content:  template.HTML(output),
	}

	err = templates.ExecuteTemplate(w, "blog.html", c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
*/

func (s *State) blogHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	bls, err := blogs_db.ListAllInfo(ctx, s.GetDB())
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

	_, err := blogs_db.Create(ctx, s.GetDB(), title, abstract, content)
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
