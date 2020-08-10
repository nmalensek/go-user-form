package main

import (
	"flag"
	"log"
	"net/http"
	"regexp"
)

var userFilePath = flag.String("ufile", "", "The absolute path for the file to use as a pseudo-database")
var validPath = regexp.MustCompile("^/(users)/([a-zA-Z0-9]+)$")

//Check the requested path; if it's valid, process it, otherwise send a 404 error.
func makeHandler(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("/users/", makeHandler(userHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
