package users

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nmalensek/go-user-form/model"
)

type mockUsers struct{}

//GetAll retrieves all saved users.
func (mu mockUsers) GetAll() ([]string, error) {
	return []string{`{"1":{"id":1,"firstName":"test2","lastName":"testLn","organization":"marketing","email":"test@email.com"}}`}, nil
}

//Create creates a new user and saves it to the "database" file.
func (mu mockUsers) Create(u model.User) bool {
	return false
}

//Edit modifies the properties of the given user based on UI input.
func (mu mockUsers) Edit(u model.User) bool {
	return false
}

//Delete finds the specified user by ID and deletes them.
func (mu mockUsers) Delete(id int) bool {
	return false
}

func TestProcessByType(t *testing.T) {
	testStore := &mockUsers{}

	req, err := http.NewRequest("GET", "/users/", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler := http.HandlerFunc(ProcessRequestByType)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

}
