package main

import (
	"html/template"
	"io/ioutil"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var validPath = regexp.MustCompile("^/(edit|save|show|delete)/([a-zA-Z0-9\\-]+)$")

type Post struct {
	Title   string
	Content string
}

func (p *Post) Slug() string {
	// code pyramid XD
	slug := strings.Replace(
		strings.Replace(
			strings.Replace(
				strings.Replace(
					strings.Replace(
						strings.Replace(
							strings.Replace(
								strings.Replace(
									strings.Replace(
										strings.Replace(
											strings.Replace(
												strings.Replace(
													strings.Replace(
														strings.ToLower(p.Title), " ", "-", -1), 
													"?", "", -1),
												"&", "", -1),
											":", "", -1),
										"!", "", -1),
									"@", "", -1),
								"#", "", -1),
							"$", "", -1),
						"%", "", -1),
					"^", "", -1),
				"*", "", -1),
			"(", "", -1),
		")", "", -1)
	return slug 
}

func (p *Post) Save() error {
	filename := "posts/" + p.Slug() + ".json"

	// serialize the data
	fstring, err := json.Marshal(p)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, fstring, 0600)

	if err != nil {
		return err
	}

	return nil
}

func DeletePost(slug string) error {
	filename := "posts/" + slug + ".json"
	err := os.Remove(filename)
	if err != nil {
		return err
	}
	return nil
}

func LoadPost(slug string) (*Post, error) {
	filename := "posts/" + slug + ".json"
	file, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}
	var p Post;
	err = json.Unmarshal(file, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func AllPosts() ([]*Post, error) {
	posts := make([]*Post, 0)
	files, err := ioutil.ReadDir("posts")
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		p, err := LoadPost(strings.TrimSuffix(f.Name(), ".json"))
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

func getSlug(r *http.Request) (string) {
	log.Printf("url %s", r.URL.Path)
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		log.Print("error: 'm' is nil")
	}
	log.Printf("slug is %s", m[2])
	return m[2]
}

func showHandler(w http.ResponseWriter, r *http.Request) {
	slug := getSlug(r)
	p, err := LoadPost(slug)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	renderTemplate(w, "show", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	slug := getSlug(r)
	p, err := LoadPost(slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	slug := getSlug(r)
	title := r.FormValue("title")
	content := r.FormValue("content")
	p := &Post{Title: title, Content: content}
	err := p.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/show/" + slug, http.StatusFound)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")
	p := &Post{Title: title, Content: content}
	err := p.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/show/" + p.Slug(), http.StatusFound)
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
	slug := getSlug(r)
	err := DeletePost(slug)

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
	log.Println("Starting server on localhost:5050/")
	log.Fatal(http.ListenAndServe(":5050", nil))
}
