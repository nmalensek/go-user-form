package main

import (
	"log"
	"net/http"
	"regexp"

	"github.com/nmalensek/go-user-form/config"
	"github.com/nmalensek/go-user-form/users"
)

var validPath = regexp.MustCompile("^/(users)/([a-zA-Z0-9]*)$")

//Check the requested path; if it's valid, process it, otherwise send a 404 error.
func makeHandler(fn func(w http.ResponseWriter, r *http.Request, e *config.Env), env *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, env)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request, e *config.Env) {
	users.ProcessRequestByType(w, r, e)
}

func main() {
	env := config.Start()

	http.HandleFunc("/users/", makeHandler(userHandler, env))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
