package main

import (
	"log"
	"net/http"

	"github.com/nmalensek/go-user-form/config"
	"github.com/nmalensek/go-user-form/users"
)

func userHandler(w http.ResponseWriter, r *http.Request, e *config.Env) {
	users.ProcessRequestByType(w, r, e)
}

func main() {
	env := config.Start()

	http.HandleFunc("/users/", makeHandler(userHandler, env))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
