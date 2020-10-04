package fileusermodel

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math"
	"sort"

	"github.com/nmalensek/go-user-form/model"
)

//FileModel constants.
const (
	databaseUnavailable = "The database is currently unavailable, please try again later."
)

//FileUserModel is an implementation of UserDataStore using the filesystem as a pseudo-database.
type FileUserModel struct {
	Filepath string
}

//GetAll retrieves all saved users.
func (m *FileUserModel) GetAll() ([]model.User, error) {
	fileData, err := readUserFile(m.Filepath)
	if err != nil {
		return nil, err
	}

	currUsers, _, err := JSONToUsers(fileData)
	if err != nil {
		return nil, err
	}

	return currUsers, nil
}

//Create creates a new user and saves it to the "database" file.
func (m *FileUserModel) Create(u *model.User) error {
	userMap, err := readFileToMap(m.Filepath)
	if err != nil {
		return err
	}

	u.ID = GetNextID(userMap)

	userMap[u.ID] = *u

	err = saveMapToFile(m.Filepath, userMap)
	if err != nil {
		return err
	}

	return nil
}

//Edit modifies the properties of the given user based on UI input.
func (m *FileUserModel) Edit(u model.User, id int) error {
	return errors.New("edit: not implemented yet")
}

//Delete finds the specified user by ID and deletes them.
func (m *FileUserModel) Delete(id int) error {
	userMap, err := readFileToMap(m.Filepath)
	if err != nil {
		return err
	}

	_, ok := userMap[id]
	if !ok {
		return errors.New(model.CouldNotFind)
	}

	delete(userMap, id)

	err = saveMapToFile(m.Filepath, userMap)
	if err != nil {
		return err
	}

	return nil
}

func readUserFile(path string) ([]byte, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return nil, errors.New(databaseUnavailable)
	}

	//copying prevents the whole file from staying in memory,
	//which is unnecessary right now because it's returning all users
	//anyway instead of a subset. Done here to get in the habit of doing this.
	cop := make([]byte, len(content))
	copy(cop, content)
	return cop, nil
}

func saveMapToFile(path string, u map[int]model.User) error {
	userBytes, err := json.Marshal(u)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, userBytes, 0644)
	if err != nil {
		log.Fatal(err)
		return errors.New(databaseUnavailable)
	}
	return nil
}

func readFileToMap(path string) (map[int]model.User, error) {
	fileData, err := readUserFile(path)
	if err != nil {
		return nil, err
	}

	_, users, err := JSONToUsers(fileData)
	if err != nil {
		return nil, err
	}
	return users, nil
}

//JSONToUsers takes a JSON string of users, puts them in a map
//keyed on ID, then sorts by ID. Returns a slice and a map of users.
func JSONToUsers(sourceBytes []byte) ([]model.User, map[int]model.User, error) {
	userMap, err := JSONToUserMap(sourceBytes)
	if err != nil {
		return nil, nil, err
	}

	userSlice := make([]model.User, len(userMap))
	i := 0
	for _, val := range userMap {
		userSlice[i] = val
		i++
	}

	sort.Slice(userSlice, func(j, k int) bool {
		return userSlice[j].ID < userSlice[k].ID
	})

	return userSlice, userMap, nil
}

//JSONToUserMap takes a JSON string of users, allocates a new map, and populates the map with the JSON string contents.
func JSONToUserMap(sourceBytes []byte) (map[int]model.User, error) {
	var userMap map[int]model.User
	err := json.Unmarshal(sourceBytes, &userMap)
	if err != nil {
		return nil, err
	}

	return userMap, nil
}

//GetNextID returns the next ID value to be assigned (current max ID + 1).
func GetNextID(userMap map[int]model.User) int {
	maxID := math.MinInt32
	for _, val := range userMap {
		if val.ID > maxID {
			maxID = val.ID
		}
	}
	return maxID + 1
}
