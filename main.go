package main

import (
	"log"
	"regexp"
	"io/ioutil"
	"net/http"
	"errors"
	"html/template"
)
//dds
var validPath = regexp.MustCompile("^/(edit|save|show)/([a-zA-Z0-9]+)$")

type Post struct {
	Title string
	Content string
}

func (p *Post) save() error {
	filename := "posts/" + p.Title + ".txt"
	err := ioutil.WriteFile(filename, []byte(p.Content), 0600)
	
	if err != nil {
		return err
	}

	return nil
}

func loadPost(title string) (*Post, error) {
	filename := "posts/" + title + ".txt"
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	p := &Post{Title: title, Content: string(content)}

	return p, nil
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
        http.NotFound(w, r)
        return "", errors.New("Invalid Post Title")
    }
    return m[2], nil // The title is the second subexpression.
}


func showHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p, err := loadPost(title)
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
	p, err := loadPost(title)
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
	err = p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/show/" + title, http.StatusFound)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")
	p := &Post{Title: title, Content: content} 
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/show/" + title, http.StatusFound)
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

func main() {
	test := &Post{Title: "test", Content: "Hello, world!"}
	err := test.save()

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/show/", showHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/new/", newHandler)
	http.HandleFunc("/create/", createHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}