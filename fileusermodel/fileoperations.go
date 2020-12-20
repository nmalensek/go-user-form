package fileusermodel

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"

	"github.com/nmalensek/go-user-form/model"
)

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
	_, users, err := fileToUsers(path)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func readFileToSlice(path string) ([]model.User, error) {
	users, _, err := fileToUsers(path)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func fileToUsers(path string) ([]model.User, map[int]model.User, error) {
	fileData, err := readUserFile(path)
	if err != nil {
		return nil, nil, err
	}

	uList, uMap, err := JSONToUsers(fileData)
	if err != nil {
		return nil, nil, err
	}
	return uList, uMap, nil
}
