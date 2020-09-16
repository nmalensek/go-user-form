package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/nmalensek/go-user-form/config"
	"github.com/nmalensek/go-user-form/users"
)

func userHandler(w http.ResponseWriter, r *http.Request, e *config.Env) {
	users.ProcessRequestByType(w, r, e)
}

func main() {
	flag.Parse()
	env, err := config.Start()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/users/", config.MakeHandler(userHandler, env))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
