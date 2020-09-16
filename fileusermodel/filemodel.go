package fileusermodel

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"sort"

	"github.com/nmalensek/go-user-form/model"
)

//FileUserModel is an implementation of UserDataStore using the filesystem as a pseudo-database.
type FileUserModel struct {
	Filepath string
}

//GetAll retrieves all saved users.
func (m *FileUserModel) GetAll() ([]model.User, error) {
	currUsers, err := readUserFile(m.Filepath)
	if err != nil {
		return nil, err
	}
	return currUsers, nil
}

//Create creates a new user and saves it to the "database" file.
func (m *FileUserModel) Create(u model.User) error {
	return errors.New("create: not implemented yet")
}

//Edit modifies the properties of the given user based on UI input.
func (m *FileUserModel) Edit(u model.User) error {
	return errors.New("edit: not implemented yet")
}

//Delete finds the specified user by ID and deletes them.
func (m *FileUserModel) Delete(id int) error {
	return errors.New("delete: not implemented yet")
}

func readUserFile(path string) ([]model.User, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var userMap map[int]model.User
	json.Unmarshal([]byte(content), &userMap)

	userSlice := make([]model.User, len(userMap))
	i := 0
	for _, val := range userMap {
		userSlice[i] = val
		i++
	}

	sort.Slice(userSlice, func(i, j int) bool {
		return userSlice[i].ID < userSlice[j].ID
	})

	return userSlice, nil
}
