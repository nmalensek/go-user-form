package users

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nmalensek/go-user-form/config"
	"github.com/nmalensek/go-user-form/model"
	"github.com/nmalensek/go-user-form/validation"
)

//ProcessRequestByType checks which HTTP verb the request has and processes it accordingly.
func ProcessRequestByType(w http.ResponseWriter, r *http.Request, e *config.Env) {
	switch r.Method {
	case http.MethodGet:
		if u, err := processGet(r, e.Datastore); err != nil {
			handleLogError(w, err, e.ErrorLog)
		} else {
			w.Write(u)
		}
	case http.MethodPost:
		if err := processPost(r, e.Datastore); err != nil {
			handleLogError(w, err, e.ErrorLog)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	case http.MethodPut:
		if err := processPut(r, e.Datastore); err != nil {
			handleLogError(w, err, e.ErrorLog)
		} else {
			w.WriteHeader(http.StatusOK)
		}
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
	user, errs := validBodyToUser(r)
	if errs != nil {
		return errs
	}

	err := db.Create(user)
	if err != nil {
		return err
	}
	return nil
}

//processPut runs validation methods, then returns nil
//if the put was successful or an error if one occurred.
func processPut(r *http.Request, db model.UserDataStore) error {
	// u, valErrs := validBodyToUser(r)
	// if valErrs != nil {
	// 	return valErrs
	// }

	////TODO: get id from query string (no key)
	// err := db.Edit(*u, )
	// if err != nil {
	// 	return err
	// }

	return nil
}

//processDelete checks for the user in the database and deletes them if
//present or returns an error if they're not found.
func processDelete(r *http.Request, db model.UserDataStore) error {

	return nil
}

func validBodyToUser(r *http.Request) (*model.User, error) {
	newUser := model.User{}
	json.NewDecoder(r.Body).Decode(&newUser)

	inputErrors := validation.ValidateInput(newUser)

	if len(inputErrors) > 0 {
		return nil, validation.UserErrors{Message: "Invalid input received, see ErrorList for details.", ErrorList: inputErrors}
	}

	return &newUser, nil
}

//handleError logs the error that occurred, writes a 500 HTTP code response header, then sends details about the error back to the requestor if applicable.
func handleLogError(w http.ResponseWriter, e error, log *log.Logger) {
	log.Println(e)

	var resp []byte
	switch e.(type) {
	case validation.UserErrors:
		data, err := json.Marshal(e)
		if err != nil {
			log.Println(err.Error())
			resp = []byte("An error occurred while processing your request, please try again later.")
		} else {
			resp = data
		}
	default:
		resp = []byte(e.Error())
	}

	w.WriteHeader(500)
	w.Write([]byte(resp))
}
