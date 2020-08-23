package model

import (
	"errors"
	"fmt"
)

//UserDataStore defines the User type data operations.
type UserDataStore interface {
	GetAll() ([]*User, error)
	Create(User) error
	Edit(User) error
	Delete(int) error
}

//User is an instance of an employee in a company.
type User struct {
	FirstName    string
	LastName     string
	Email        string
	Organization string
}

func (u User) String() string {
	return fmt.Sprintf("%v\t%v\t%v\t%v", u.FirstName, u.LastName, u.Email, u.Organization)
}

//FileUserModel is an implementation of UserDataStore using the filesystem as a pseudo-database.
type FileUserModel struct {
	Filepath *string
}

//GetAll retrieves all saved users.
func (m *FileUserModel) GetAll() ([]*User, error) {
	//content, err := ioutil.ReadFile(m.Filepath)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return nil, err
	// }
	return []*User{}, nil
}

//Create creates a new user and saves it to the "database" file.
func (m *FileUserModel) Create(u User) error {
	return errors.New("create: not implemented yet")
}

//Edit modifies the properties of the given user based on UI input.
func (m *FileUserModel) Edit(u User) error {
	return errors.New("edit: not implemented yet")
}

//Delete finds the specified user by ID and deletes them.
func (m *FileUserModel) Delete(id int) error {
	return errors.New("delete: not implemented yet")
}
