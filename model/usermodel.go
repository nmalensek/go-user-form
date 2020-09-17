package model

import (
	"fmt"
)

//UserDataStore defines the User type data operations.
type UserDataStore interface {
	GetAll() ([]User, error)
	Create(*User) error
	Edit(User) error
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
	return fmt.Sprintf("%v\t%v\t%v\t%v", u.FirstName, u.LastName, u.Email, u.Organization)
}
