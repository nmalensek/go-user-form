package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nmalensek/go-user-form/fileusermodel"

	"github.com/nmalensek/go-user-form/config"
	"github.com/nmalensek/go-user-form/model"
)

const (
	baseMockData = `{"1":{"id":1,"firstName":"test","lastName":"testLn","organization":"marketing","email":"test@email.com"},"2":{"id":2,"firstName":"test2","lastName":"testLn","organization":"sales","email":"new@employee.com"}}`
)

type mockUsers struct {
	UserData string
}

//GetAll retrieves all saved users.
func (mu *mockUsers) GetAll() ([]model.User, error) {
	mockJSON := mu.dataSet()

	userSlice, _, err := fileusermodel.JSONToUsers([]byte(mockJSON))
	if err != nil {
		return nil, err
	}

	return userSlice, nil
}

//Create creates a new user and saves it to the "database" file.
func (mu *mockUsers) Create(u *model.User) error {
	mockData := mu.dataSet()

	_, userMap, err := fileusermodel.JSONToUsers([]byte(mockData))
	if err != nil {
		return err
	}

	var newID = fileusermodel.GetNextID(userMap)

	if newID <= 0 {
		return errors.New("error in create: user ID should not be less than or equal to 0")
	}

	u.ID = newID

	JSONUser, err := u.JSONString()
	if err != nil {
		return err
	}

	mu.setDataSet(string(JSONUser), true)

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
		mu.UserData = baseMockData
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
	mockEnv := config.Env{Datastore: testStore}

	req, err := http.NewRequest(http.MethodPost, "/users/",
		strings.NewReader(`{"firstName":"testUser","lastName":"test1","email":"test@email.com","organization":"sales"}`))
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

	want := model.User{ID: 3, FirstName: "testUser", LastName: "test1",
		Email: "test@email.com", Organization: "sales"}
	got := userMap[3]
	if want != got {
		t.Errorf("incorrect user save data, got %v want %v",
			got, want)
	}

}

func TestPostProcessingInvalid(t *testing.T) {
	//test server-side validation.
}
