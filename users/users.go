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
	case http.MethodGet:
		if u, err := processGet(r, e.Datastore); err != nil {
			e.ErrorLog.Println(err)
			handleError(w, err)
		} else {
			w.Write(u)
		}
	case http.MethodPost:
		if err := processPost(r, e.Datastore); err != nil {
			handleError(w, err)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	case http.MethodPut:
		//TODO
	case http.MethodDelete:
		//TODO
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

//processPost runs validation methods, then returns nil
//if the post was successful or an error if one occurred.
func processPost(r *http.Request, db model.UserDataStore) error {
	newUser := model.User{}
	json.NewDecoder(r.Body).Decode(&newUser)
	//TODO: validate.
	err := db.Create(&newUser)
	if err != nil {
		return err
	}
	return nil
}

//processPut runs validation methods, then returns nil
//if the put was successful or an error if one occurred.
func processPut(r *http.Request, db model.UserDataStore) error {
	return nil
}

//TODO: jsonify errors because the front end needs them that way.
func handleError(w http.ResponseWriter, e error) {
	w.WriteHeader(500)
	w.Write([]byte(e.Error()))
}
