package main

import (
	"websites/sites"
	"net/http"
	"fmt"
	"strings"
)

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Received host: "+r.Host)
    fmt.Println("Received path: "+r.URL.Path)

		title := string("/")
		if len(r.URL.Path) > 2 {
			path := strings.TrimRight(r.URL.Path, "/")
			last := path[strings.LastIndex(path,"/")+1:]
			title = last
			fmt.Println("title: "+title);
		}
//	if !titleValidator.MatchString(title) {
//			http.NotFound(w, r)
//			return
//		}
		fn(w, r, title)
	}
}

func main() {
  fmt.Println("Registering handlers");
	http.HandleFunc("makebook.41j.com/down/", makeHandler(sites.MakebookDownHandler))
	http.HandleFunc("makebook.41j.com/form/", makeHandler(sites.MakebookFormHandler))
	http.HandleFunc("makebook.41j.com/"     , makeHandler(sites.MakebookRootHandler))
	http.HandleFunc("graph.41j.com/down/"   , makeHandler(sites.GraphJSONDownHandler))
	http.HandleFunc("graph.41j.com/graph/"  , makeHandler(sites.GraphJSONFormHandler))
	http.HandleFunc("graph.41j.com/"        , makeHandler(sites.GraphJSONRootHandler))
	http.ListenAndServe(":80", nil)
}
