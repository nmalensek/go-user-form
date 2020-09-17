package users

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nmalensek/go-user-form/fileusermodel"

	"github.com/nmalensek/go-user-form/config"
	"github.com/nmalensek/go-user-form/model"
)

type mockUsers struct {
	UserData string
}

//GetAll retrieves all saved users.
func (mu *mockUsers) GetAll() ([]model.User, error) {
	mockJSON := mu.dataSet()

	var userMap map[int]model.User

	userSlice, err := fileusermodel.JSONToUsers([]byte(mockJSON), &userMap)
	if err != nil {
		return nil, err
	}

	return userSlice, nil
}

//Create creates a new user and saves it to the "database" file.
func (mu *mockUsers) Create(u model.User) error {
	return nil
}

//Edit modifies the properties of the given user based on UI input.
func (mu *mockUsers) Edit(u model.User) error {
	return nil
}

//Delete finds the specified user by ID and deletes them.
func (mu *mockUsers) Delete(id int) error {
	return nil
}

func (mu *mockUsers) dataSet() string {
	if mu.UserData == "" {
		mu.UserData = `{"1":{"id":1,"firstName":"test","lastName":"testLn","organization":"marketing","email":"test@email.com"},"2":{"id":2,"firstName":"test2","lastName":"testLn","organization":"sales","email":"new@employee.com"}}`
	}
	return mu.UserData
}

func (mu *mockUsers) setDataSet(newData string, append bool) {
	if append {
		var userMap map[int]model.User
		json.Unmarshal([]byte(mu.dataSet()), &userMap)
		var newUser model.User
		json.Unmarshal([]byte(newData), &newUser)

		userMap[newUser.ID] = newUser
		updatedData, err := json.Marshal(userMap)
		if err != nil {
			return
		}
		mu.UserData = string(updatedData)
	} else {
		mu.UserData = newData
	}
}

//Test getting all users; method should return all users in the datastore.
func TestGetProcessing(t *testing.T) {
	testStore := &mockUsers{}
	mockEnv := config.Env{Datastore: testStore}

	req, err := http.NewRequest(http.MethodGet, "/users/", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler := http.HandlerFunc(config.MakeHandler(ProcessRequestByType, &mockEnv))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rec.Code, http.StatusOK)
	}

	want := `[{"id":1,"firstName":"test","lastName":"testLn","organization":"marketing","email":"test@email.com"},{"id":2,"firstName":"test2","lastName":"testLn","organization":"sales","email":"new@employee.com"}]`
	got := rec.Body.String()

	if strings.EqualFold(got, want) {
		t.Errorf("handler returned wrong content, got %v want %v",
			got, want)
	}
}

//Test new user creation; new user should be added to the datastore.
func TestPostProcessingGood(t *testing.T) {
	testStore := &mockUsers{}
	//make sure the test store is empty for ease of save checking.
	testStore.setDataSet("", false)
	mockEnv := config.Env{Datastore: testStore}

	req, err := http.NewRequest(http.MethodPost, "/users/",
		strings.NewReader(`{"firstName":"testUser","lastName":"test1","email":"test@email1.com","organization":"sales"}`))
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(config.MakeHandler(ProcessRequestByType, &mockEnv))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rec.Code, http.StatusOK)
	}

	updatedData := testStore.dataSet()

	var userMap map[int]model.User
	json.Unmarshal([]byte(updatedData), &userMap)

	want := model.User{FirstName: "testUser", LastName: "test1",
		Email: "test@email.com", Organization: "sales"}
	got := userMap[1]
	if want != got {
		t.Errorf("problem during user save, got %v want %v",
			got, want)
	}

}

func TestPostProcessingInvalid(t *testing.T) {
	//test server-side validation.
}
