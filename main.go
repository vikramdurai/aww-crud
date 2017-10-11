/*
blog-app-go is a simple CRUD (Create-Read-Update-Delete) application
using the net/http library.
*/
package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var validPath = regexp.MustCompile("^/(edit|save|show|delete)/([a-zA-Z0-9]+)$")

// Post is how posts are internally stored
type Post struct {
	Title   string
	Content string
}

// Save writes the Post to the "posts"
// directory
func (p *Post) Save() error {
	filename := "posts/" + p.Title + ".txt"
	err := ioutil.WriteFile(filename, []byte(p.Content), 0600)

	if err != nil {
		return err
	}

	return nil
}

// DeletePost deletes a Post using the given "title"
// from the "posts" directory
func DeletePost(title string) error {
	filename := "posts/" + title + ".txt"
	err := os.Remove(filename)
	if err != nil {
		return err
	}
	return nil
}

// loadPost loads a Post using the given "title"
// from the "posts" directory
func LoadPost(title string) (*Post, error) {
	filename := "posts/" + title + ".txt"
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	p := &Post{Title: title, Content: string(content)}

	return p, nil
}

// AllPosts returns an array of all posts
// by looking into the "posts" directory
func AllPosts() ([]*Post, error) {
	posts := make([]*Post, 1)
	files, err := ioutil.ReadDir("posts")
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		p, err := LoadPost(strings.Trim(f.Name(), ".txt"))
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}
	return posts, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Post) {
	t, err := template.ParseFiles("templates/" + tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		log.Fatal("this should never happen")
	}
	return m[2], nil
}

func showHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p, err := LoadPost(title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "show", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p, err := LoadPost(title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	content := r.FormValue("content")
	p := &Post{Title: title, Content: content}
	err = p.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/show/"+title, http.StatusFound)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")
	p := &Post{Title: title, Content: content}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/show/"+title, http.StatusFound)
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/new.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = DeletePost(title)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := AllPosts()

	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		defer log.Fatal(err)
		return
	}

	t, err := template.ParseFiles("templates/index.html")

	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		// return
		defer log.Fatal(err)
		return
	}

	err = t.Execute(w, posts)

	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		defer log.Fatal(err)
		return
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/show/", showHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/new/", newHandler)
	http.HandleFunc("/create/", createHandler)
	http.HandleFunc("/delete/", deleteHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
