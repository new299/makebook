package sites

import (
	"html/template"
	"net/http"
  "fmt"
	"crypto/md5"
	"hash"
    "os"
	"os/exec"
	"strings"
	"io/ioutil"
	"encoding/hex"
)

func GraphJSONDownHandler(w http.ResponseWriter, r *http.Request, title string) {
	graphdata := r.FormValue("graphdata")
	
	var h hash.Hash = md5.New()
	h.Write([]byte(graphdata))
	
  filename :=  hex.EncodeToString(h.Sum(nil))
 
  fmt.Println("writing graph: " + filename)
  f, err := os.Create(filename)
  if err != nil {
    fmt.Println(err)
  }

	data := strings.Split(graphdata," ")
  for n:=0;n<len(data);n++ {
		fmt.Fprintf(f,"%d %s\n",n,data[n]);
	}

  f.Close()

  fmt.Println("writing gnuplot: " + filename + ".gnuplot");
  gf, err := os.Create(filename + ".gnuplot")
  if err != nil {
    fmt.Println(err)
  }

  fmt.Fprintf(gf,"set terminal png\n");
	fmt.Fprintf(gf,"set nokey\n");
	fmt.Fprintf(gf,"set output \"%s\"\n",filename + ".png");
	fmt.Fprintf(gf,"plot \"%s\" using 1:2 with lines",filename);
  gf.Close()

  cmd := exec.Command("gnuplot",filename + ".gnuplot")
  out,e := cmd.Output()
  if e != nil { fmt.Println(e); fmt.Println(out); }


  w.Header().Set("Content-Type", "image/png")

  graphimg, e := ioutil.ReadFile(filename + ".png")
  if e != nil {
    fmt.Println("read error: " + filename + ".png")
    return
  }

  w.Write(graphimg);


}

func GraphJSONRootHandler(w http.ResponseWriter, r *http.Request, title string) {
	fmt.Println("graphroothandler");
	http.Redirect(w, r, "/graph/", http.StatusFound)
}

func GraphJSONFormHandler(w http.ResponseWriter, r *http.Request, title string) {

  fmt.Println("graphhandler");
	p, err := loadPage("graph")
	if err != nil {
		p = &Page{Title: title}
	}
	renderGraphTemplate(w, "graph", p)
}

var graphtemplates = template.Must(template.ParseFiles("html/graph.html"))

func renderGraphTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := graphtemplates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
