/*
Presenting aww-crud: my first CRUD[1] app in go

[1]: Full form is Create, Update, Delete
*/

package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var validPath = regexp.MustCompile("^/(edit|save|show|delete)/([a-zA-Z0-9\\-]+)$")

type Record struct {
	Title   string
	Content string
}

func (r *Record) Slug() string {
	slug := strings.ToLower(r.Title)

	// replace spaces with -
	slug = strings.Replace(slug, " ", "-", -1)

	// strip unwanted characters
	re := regexp.MustCompile("[?&:!@#$%^*()]")

	return re.ReplaceAllLiteralString(slug, "")
}

func (r *Record) Save() error {
	filename := "records/" + r.Slug() + ".json"

	// serialize the data
	fstring, err := json.Marshal(r)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, fstring, 0600)

	if err != nil {
		return err
	}

	return nil
}

func DeleteRecord(slug string) error {
	filename := "records/" + slug + ".json"
	err := os.Remove(filename)
	if err != nil {
		return err
	}
	return nil
}

func LoadRecord(slug string) (*Record, error) {
	filename := "records/" + slug + ".json"
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var r Record
	err = json.Unmarshal(file, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func AllRecords() ([]*Record, error) {
	records := make([]*Record, 0)
	files, err := ioutil.ReadDir("records")
	if os.IsNotExist(err) {
		if err := os.Mkdir("records", os.ModePerm); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	for _, f := range files {
		r, err := LoadRecord(strings.TrimSuffix(f.Name(), ".json"))
		if err != nil {
			return nil, err
		}

		records = append(records, r)
	}
	return records, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, r *Record) {
	t, err := template.ParseFiles("templates/" + tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getSlug(r *http.Request) string {
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
	rec, err := LoadRecord(slug)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	renderTemplate(w, "show", rec)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	slug := getSlug(r)
	rec, err := LoadRecord(slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "edit", rec)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")
	rec := &Record{Title: title, Content: content}
	err := rec.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/show/"+rec.Slug(), http.StatusFound)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")
	rec := &Record{Title: title, Content: content}
	err := rec.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/show/"+rec.Slug(), http.StatusFound)
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
	err := DeleteRecord(slug)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	records, err := AllRecords()
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalf("unable to load all records: %v", err)
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		// return
		log.Fatalf("unable to parse file: %v", err)
	}

	err = t.Execute(w, records)

	if err != nil {
		log.Fatalf("unable to render template: %v", err)
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
