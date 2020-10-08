package model

import (
	"encoding/json"
	"fmt"
)

//Model error message constants.
const (
	CouldNotFind          = "Could not find specified user in database."
	CreateErrorBadID      = "Could not create user, unable to assign valid ID."
	CreateErrorIncomplete = "Could not create user from the information provided."
)

//UserDataStore defines the User type data operations.
type UserDataStore interface {
	GetAll() ([]User, error)
	Create(*User) error
	Edit(User, int) error
	Delete(int) error
}

//User is an instance of an employee in a company.
type User struct {
	ID           int    `json:"id"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	Organization string `json:"organization"`
}

func (u User) String() string {
	return fmt.Sprintf("%v\t%v\t%v\t%v\t%v", u.ID, u.FirstName, u.LastName, u.Email, u.Organization)
}

//JSONString returns the JSON version of the user.
func (u User) JSONString() ([]byte, error) {
	JSONUser, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	return JSONUser, nil
}
