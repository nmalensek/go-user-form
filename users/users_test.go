package users

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nmalensek/go-user-form/config"
	"github.com/nmalensek/go-user-form/model"
)

type mockUsers struct{}

//GetAll retrieves all saved users.
func (mu mockUsers) GetAll() ([]model.User, error) {
	mockJSON := `{"1":{"id":1,"firstName":"test","lastName":"testLn","organization":"marketing","email":"test@email.com"},"2":{"id":2,"firstName":"test2","lastName":"testLn","organization":"sales","email":"new@employee.com"}}`
	var userArr []model.User
	json.Unmarshal([]byte(mockJSON), &userArr)

	return userArr, nil
}

//Create creates a new user and saves it to the "database" file.
func (mu mockUsers) Create(u model.User) error {
	return nil
}

//Edit modifies the properties of the given user based on UI input.
func (mu mockUsers) Edit(u model.User) error {
	return nil
}

//Delete finds the specified user by ID and deletes them.
func (mu mockUsers) Delete(id int) error {
	return nil
}

func TestProcessByType(t *testing.T) {
	testStore := &mockUsers{}
	mockEnv := config.Env{Datastore: testStore}

	req, err := http.NewRequest("GET", "/users/", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler := http.HandlerFunc(config.MakeHandler(ProcessRequestByType, &mockEnv))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

}
