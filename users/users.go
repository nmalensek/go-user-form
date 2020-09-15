package users

import (
	"encoding/json"
	"net/http"

	"github.com/nmalensek/go-user-form/config"
	"github.com/nmalensek/go-user-form/model"
)

//ProcessRequestByType checks which HTTP verb the request has and processes it accordingly.
func ProcessRequestByType(w http.ResponseWriter, r *http.Request, e *config.Env) {
	switch r.Method {
	case "GET":
		if u, err := processGet(r, e.Datastore); err != nil {
			e.ErrorLog.Println(err)
			handleError(w, err)
		} else {
			w.Write(u)
		}
	}
}

//processGet returns bytes from JSON records from the database or an error if one occurs.
func processGet(r *http.Request, db model.UserDataStore) ([]byte, error) {
	userList, err := db.GetAll()
	if err != nil {
		return nil, err
	}

	userBytes, err := json.Marshal(userList)
	if err != nil {
		return nil, err
	}
	return userBytes, nil
}

//TODO: jsonify errors because the front end needs them that way.
func handleError(w http.ResponseWriter, e error) {
	w.WriteHeader(500)
	w.Write([]byte(e.Error()))
}
