package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"fmt"
	"os"
	"io"
)

type Page struct {
	Title string
	Body  []byte
}

func makebook(bookname string) error {
	filename := bookname + ".epub"

  fmt.Println("writing: " + filename);
	f, err := os.Create(filename)
  if err != nil { fmt.Println(err); return err }
  n, err := io.WriteString(f, "blahblahblah")
  if err != nil { fmt.Println(n, err); return err }
  f.Close()

  return err
}

func loadPage(title string) (*Page, error) {
	filename := title
	body, err := ioutil.ReadFile(filename+".epub")

	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func downHandler(w http.ResponseWriter, r *http.Request, title string) {

  booktitle := r.FormValue("booktitle")

  makebook(booktitle)

  bookdata, e := ioutil.ReadFile(booktitle + ".epub")
  if e != nil {
    fmt.Println("read error: " + booktitle)
  }

  w.Header().Set("Content-Type", "application/xhtml+xml")

	w.Write(bookdata);
}

func formHandler(w http.ResponseWriter, r *http.Request, title string) {

  fmt.Println("formhandler");
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "form", p)
}

var templates = template.Must(template.ParseFiles("form.html", "down.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const lenPath = len("/down/")

var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

    fmt.Println("Received path: "+r.URL.Path);
		title := r.URL.Path[lenPath:]
	//	if !titleValidator.MatchString(title) {
//			http.NotFound(w, r)
//			return
//		}
		fn(w, r, title)
	}
}

func main() {
  fmt.Println("Registering handlers");
	http.HandleFunc("/down/", makeHandler(downHandler))
	http.HandleFunc("/form/", makeHandler(formHandler))
	http.ListenAndServe(":80", nil)
}
