package fileusermodel

import (
	"errors"
	"math"

	"github.com/nmalensek/go-user-form/validation"

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
	users, err := readFileToSlice(m.Filepath)
	if err != nil {
		return nil, err
	}
	return users, nil
}

//Create creates a new user and saves it to the "database" file.
func (m *FileUserModel) Create(u *model.User) error {
	errs := validation.ValidateCompleteInput(*u)
	if len(errs) > 0 {
		return errors.New(model.CreateErrorIncomplete)
	}

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
	errs := validation.ValidatePartialInput(u)
	if len(errs) > 0 {
		return errors.New(model.EditErrorIncomplete)
	}

	userMap, err := readFileToMap(m.Filepath)
	if err != nil {
		return err
	}

	_, ok := userMap[id]
	if !ok {
		return errors.New(model.CouldNotFind)
	}

	savedUser := userMap[id]
	if u.FirstName != "" {
		savedUser.FirstName = u.FirstName
	}
	if u.LastName != "" {
		savedUser.LastName = u.LastName
	}
	if u.Email != "" {
		savedUser.Email = u.Email
	}
	if u.Organization != "" {
		savedUser.Organization = u.Organization
	}

	userMap[id] = savedUser

	err = saveMapToFile(m.Filepath, userMap)
	if err != nil {
		return err
	}

	return nil
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
