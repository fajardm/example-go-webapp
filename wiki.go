package main

import (
	"io/ioutil"
	"net/http"
	"html/template"
)

type Page struct {
	Title string
	Body  []byte
}

func (page *Page) save() error {
	filename := page.Title + ".txt"
	return ioutil.WriteFile(filename, page.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(res http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len("/view/"):]
	page, err := loadPage(title)
	if err != nil {
		http.Redirect(res, req, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(res, "view", page)
}

func editHandler(res http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len("/edit/"):]
	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}
	renderTemplate(res, "edit", page)
}

func saveHandler(res http.ResponseWriter, req *http.Request) {
	title := req.URL.Path[len("/save/"):]
	body := req.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(res, req, "/view/"+title, http.StatusFound)
}

func renderTemplate(res http.ResponseWriter, v string, page *Page) {
	view, _ := template.ParseFiles(v + ".html")
	view.Execute(res, page)
}

func main() {
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)
}
