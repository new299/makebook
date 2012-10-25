package sites

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"fmt"
	"os"
	"io"
	"os/exec"
	"strings"
)

type Page struct {
	Title string
	Body  []byte
}

func makebook(bookname string) error {

	search := strings.Replace(bookname, " ", "%20", -1)

  cmdargs := "http://en.wikipedia.org/w/index.php?title=Special%3ASearch&profile=default&search=" + search + "&fulltext=Search"
  cmd := exec.Command("curl",cmdargs)
  out,e := cmd.Output()
	if e != nil { fmt.Println(e) }

	s := string(out);
  re := regexp.MustCompile(`mw-search-result-heading'><a href="[A-Za-z0-9_//]+"`)
  searchres := re.FindAllStringSubmatch(s,10)

 	trimlen := len("mw-search-result-heading'><a href=\"/wiki/")

  exec.Command("mkdir",bookname).Output()

  fmt.Println("writing: " + bookname + ".html");
	indexfile, err := os.Create(bookname + ".html")
  if err != nil {fmt.Println(err);}

  for i := 0; i < len(searchres); i++ {
		current := string(searchres[i][0])

    pagename := current[trimlen:len(current)-1]
  	url := "http://en.wikipedia.org/w/index.php?title=" + current[trimlen:len(current)-1] + "&printable=yes"

    fmt.Println(url);
    wgetcmd := exec.Command("wget","--wait=1","-k","-E","-H","-K","-e","robots=off","-p",url)
		wgetcmd.Dir = bookname
		wgetout,e := wgetcmd.CombinedOutput()
		fmt.Println(string(wgetout))
    if e != nil { fmt.Println(e); }

    io.WriteString(indexfile,"<a href=\"");
    io.WriteString(indexfile,bookname + "/en.wikipedia.org/w/index.php?title=" + pagename + "&printable=yes.html");
		io.WriteString(indexfile,"\">" + pagename + "</a><br>");
  }

  indexfile.Close()

  ebookcmd := exec.Command("ebook-convert",bookname + ".html",bookname + ".epub")
	ebookcmd.Output()

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

func MakebookDownHandler(w http.ResponseWriter, r *http.Request, title string) {

  booktitle := r.FormValue("booktitle")
	w.Header().Set("Content-Type", "application/epub+zip")

	w.Header().Set("Content-Disposition", "attachment; filename=\"" + booktitle +".epub\"")

  bookdata, e := ioutil.ReadFile(booktitle + ".epub")
  if e != nil {
    fmt.Println("read error: " + booktitle)
  	makebook(booktitle)
  	bookdata, e := ioutil.ReadFile(booktitle + ".epub")
		if e != nil { fmt.Println(e) }

		w.Write(bookdata);
		return
  }

	w.Write(bookdata);
}

func MakebookRootHandler(w http.ResponseWriter, r *http.Request, title string) {
	http.Redirect(w, r, "/form/", http.StatusFound)
}

func MakebookFormHandler(w http.ResponseWriter, r *http.Request, title string) {

  fmt.Println("formhandler");
	p, err := loadPage("form")
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "form", p)
}

var templates = template.Must(template.ParseFiles("form.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
