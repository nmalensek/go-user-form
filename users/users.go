package users

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"

	"github.com/nmalensek/go-user-form/config"
	"github.com/nmalensek/go-user-form/model"
	"github.com/nmalensek/go-user-form/validation"
)

//Handler error messages.
const (
	MalformedURI         = "Received malformed URI, please check input and try again"
	InvalidInput         = "Invalid input received, see ErrorList for details."
	ErrorWhileProcessing = "An error occurred while processing your request, please try again later."
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
		if err := processDelete(r, e.Datastore); err != nil {
			handleLogError(w, err, e.ErrorLog)
		} else {
			w.WriteHeader(http.StatusOK)
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
	u, valErrs := validBodyToUser(r)
	if valErrs != nil {
		return valErrs
	}

	id, ok := getIDFromPath(r.URL.EscapedPath())
	if !ok {
		return errors.New(MalformedURI)
	}

	err := db.Edit(*u, id)
	if err != nil {
		return err
	}

	return nil
}

//processDelete checks for the user in the database and deletes them if
//present or returns an error if they're not found.
func processDelete(r *http.Request, db model.UserDataStore) error {
	id, ok := getIDFromPath(r.URL.EscapedPath())

	if !ok {
		return errors.New(MalformedURI)
	}

	err := db.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func validBodyToUser(r *http.Request) (*model.User, error) {
	newUser := model.User{}
	json.NewDecoder(r.Body).Decode(&newUser)

	inputErrors := validation.ValidateInput(newUser)

	if len(inputErrors) > 0 {
		return nil, validation.UserErrors{Message: InvalidInput, ErrorList: inputErrors}
	}

	return &newUser, nil
}

func getIDFromPath(p string) (int, bool) {
	//should end with after /number, don't care what comes before.
	numberPatt := regexp.MustCompile(`/([0-9]+)$`)

	//edit URI should be /users/id, so this should find the whole string and the ID if valid.
	id := numberPatt.FindStringSubmatch(p)
	if id == nil || len(id) != 2 {
		return math.MinInt32, false
	}

	parsedID, err := strconv.Atoi(id[1])
	if err != nil {
		return math.MinInt32, false
	}

	return parsedID, true
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
			resp = []byte(ErrorWhileProcessing)
		} else {
			resp = data
		}
	default:
		resp = []byte(e.Error())
	}

	w.WriteHeader(500)
	w.Write([]byte(resp))
}
