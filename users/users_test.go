package users

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nmalensek/go-user-form/fileusermodel"
	"github.com/nmalensek/go-user-form/validation"

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
		return errors.New(model.CreateErrorBadID)
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
func (mu *mockUsers) Edit(u model.User, id int) error {
	mockData := mu.dataSet()

	userMap, err := fileusermodel.JSONToUserMap([]byte(mockData))
	if err != nil {
		return err
	}

	origUser, ok := userMap[id]
	if !ok {
		return errors.New(model.CouldNotFind)
	}

	origUser.FirstName = u.FirstName
	origUser.LastName = u.LastName
	origUser.Email = u.Email
	origUser.Organization = u.Organization

	userMap[id] = origUser

	JSONBytes, err := json.Marshal(userMap)
	if err != nil {
		return err
	}

	mu.setDataSet(string(JSONBytes), false)

	return nil
}

//Delete finds the specified user by ID and deletes them.
func (mu *mockUsers) Delete(id int) error {
	mockData := mu.dataSet()

	userMap, err := fileusermodel.JSONToUserMap([]byte(mockData))
	if err != nil {
		return err
	}
	_, ok := userMap[id]
	if !ok {
		return errors.New(model.CouldNotFind)
	}

	delete(userMap, id)

	JSONBytes, _ := json.Marshal(userMap)

	mu.setDataSet(string(JSONBytes), false)

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

func makeMockEnv() config.Env {
	testStore := &mockUsers{}
	var buf bytes.Buffer
	testLogger := log.New(&buf, "test log: ", log.Lshortfile)
	return config.Env{Datastore: testStore, ErrorLog: testLogger}
}

//Test getting all users; method should return all users in the datastore.
func TestGetProcessing(t *testing.T) {
	mockEnv := makeMockEnv()

	req, err := http.NewRequest(http.MethodGet, "/users/", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler := http.HandlerFunc(config.MakeHandler(ProcessRequestByType, &mockEnv))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	compareStatusCode(rec.Code, http.StatusOK, t)

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

	compareStatusCode(rec.Code, http.StatusOK, t)

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

func TestPostMissingFields(t *testing.T) {
	mockEnv := makeMockEnv()

	req, err := http.NewRequest(http.MethodPost, "/users/",
		strings.NewReader(`{"firstName":"testUser","email":"test@email.com","organization":"sales"}`))
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(config.MakeHandler(ProcessRequestByType, &mockEnv))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	compareStatusCode(rec.Code, http.StatusInternalServerError, t)

	var errs validation.UserErrors
	json.NewDecoder(rec.Body).Decode(&errs)

	if len(errs.ErrorList) != 1 {
		t.Errorf("Expected one error, got %v errors.", len(errs.ErrorList))
	}

	want := validation.UserError{PropName: "LastName", PropValue: ""}
	got := validation.UserError{PropName: errs.ErrorList[0].PropName, PropValue: errs.ErrorList[0].PropValue}

	compareGotWant(got, want, t)
}

func TestPostInvalidEmail(t *testing.T) {
	mockEnv := makeMockEnv()

	req, err := http.NewRequest(http.MethodPost, "/users/",
		strings.NewReader(`{"firstName":"testUser", "lastName":"test","email":"test@.net","organization":"sales"}`))
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(config.MakeHandler(ProcessRequestByType, &mockEnv))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	compareStatusCode(rec.Code, http.StatusInternalServerError, t)

	var errs validation.UserErrors
	json.NewDecoder(rec.Body).Decode(&errs)

	if len(errs.ErrorList) != 1 {
		t.Errorf("Expected one error, got %v errors.", len(errs.ErrorList))
	}

	want := validation.UserError{PropName: "Email", Message: validation.IncorrectFormatMessage("Email")}
	got := validation.UserError{PropName: errs.ErrorList[0].PropName, Message: errs.ErrorList[0].Message}

	compareGotWant(got, want, t)
}

func TestPutValid(t *testing.T) {
	mockEnv := makeMockEnv()

	req, err := http.NewRequest(http.MethodPut, "/users/1",
		strings.NewReader(`{"firstName":"editTestUser", "lastName":"test","email":"test@t.net","organization":"sales"}`))
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(config.MakeHandler(ProcessRequestByType, &mockEnv))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	compareStatusCode(rec.Code, http.StatusOK, t)

	updatedData, _ := mockEnv.Datastore.GetAll()
	var got model.User
	for _, item := range updatedData {
		if item.ID == 1 {
			got = item
			break
		}
	}
	want := model.User{ID: 1, FirstName: "editTestUser", LastName: "test", Email: "test@t.net", Organization: "sales"}

	compareGotWant(got, want, t)
}

func TestPutInvalidID(t *testing.T) {
	mockEnv := makeMockEnv()

	req, err := http.NewRequest(http.MethodPut, "/users/1asdf",
		strings.NewReader(`{"firstName":"editTestUser", "lastName":"test","email":"test@t.net","organization":"sales"}`))
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(config.MakeHandler(ProcessRequestByType, &mockEnv))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	compareStatusCode(rec.Code, http.StatusInternalServerError, t)

	compareGotWant(rec.Body.String(), MalformedURI, t)
}

func TestPutMissingID(t *testing.T) {
	mockEnv := makeMockEnv()
	fakeID := math.MaxInt64
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/users/%v", fakeID),
		strings.NewReader(`{"firstName":"editTestUser", "lastName":"test","email":"test@t.net","organization":"sales"}`))
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(config.MakeHandler(ProcessRequestByType, &mockEnv))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	compareStatusCode(rec.Code, http.StatusInternalServerError, t)

	compareGotWant(rec.Body.String(), model.CouldNotFind, t)
}

func TestPutInvalidUser(t *testing.T) {
	mockEnv := makeMockEnv()

	req, err := http.NewRequest(http.MethodPut, "/users/1",
		strings.NewReader(`{"firstName":"", "email":"","organization":""}`))
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(config.MakeHandler(ProcessRequestByType, &mockEnv))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	compareStatusCode(rec.Code, http.StatusInternalServerError, t)

	var errs validation.UserErrors
	json.NewDecoder(rec.Body).Decode(&errs)

	if len(errs.ErrorList) != 1 {
		t.Errorf("Expected one error, got %v errors.", len(errs.ErrorList))
	}

	compareGotWant(errs.ErrorList[0].Message, validation.MissingAllProps, t)
}

func TestDeleteValid(t *testing.T) {
	mockEnv := makeMockEnv()

	//make sure user to delete's there to make sure delete's actually modifying collection.
	users, _ := mockEnv.Datastore.GetAll()
	var ok bool
	for _, v := range users {
		if v.ID == 1 {
			ok = true
			break
		}
	}
	if !ok {
		t.Errorf("test user not in database originally.")
	}

	req, err := http.NewRequest(http.MethodDelete, "/users/1", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(config.MakeHandler(ProcessRequestByType, &mockEnv))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	compareStatusCode(rec.Code, http.StatusOK, t)

	//user from before shouldn't be there anymore.
	users, _ = mockEnv.Datastore.GetAll()
	for _, v := range users {
		if v.ID == 1 {
			ok = false
			break
		}
	}
	compareGotWant(ok, true, t)
}

func TestDeleteInvalid(t *testing.T) {
	mockEnv := makeMockEnv()
	fakeID := math.MaxInt64
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%v", fakeID), strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(config.MakeHandler(ProcessRequestByType, &mockEnv))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	compareStatusCode(rec.Code, http.StatusInternalServerError, t)

	gotErr := rec.Body.String()
	wantErr := model.CouldNotFind

	compareGotWant(gotErr, wantErr, t)
}

func compareGotWant(got interface{}, want interface{}, t *testing.T) {
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func compareStatusCode(got int, want int, t *testing.T) {
	if got != want {
		t.Errorf("handler returned wrong status code: got %v want %v", got, want)
	}
}
