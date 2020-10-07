package fileusermodel

import (
	"io/ioutil"
	"math"
	"os"
	"testing"

	"github.com/nmalensek/go-user-form/model"
)

const (
	baseMockData = `{"1":{"id":1,"firstName":"test","lastName":"testLn","organization":"marketing","email":"test@email.com"},"2":{"id":2,"firstName":"test2","lastName":"testLn","organization":"sales","email":"new@employee.com"}}`
	testFilePath = "./testUserStore.json"
)

func TestMain(m *testing.M) {
	ioutil.WriteFile(testFilePath, []byte(baseMockData), 0644)
	m.Run()
	os.Remove(testFilePath)
}

func TestFileUnavailable(t *testing.T) {

}

func TestGetAll(t *testing.T) {
	model := FileUserModel{Filepath: testFilePath}
	mockUsers, err := model.GetAll()
	if err != nil {
		t.Errorf(err.Error())
	}

	convertedConst, _, err := JSONToUsers([]byte(baseMockData))
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(mockUsers) != len(convertedConst) {
		t.Errorf("got length %v want length %v", len(mockUsers), len(convertedConst))
	}

	for i := range mockUsers {
		if mockUsers[i] != convertedConst[i] {
			t.Errorf("expected %v, got %v", convertedConst[i], mockUsers[i])
		}
	}
}

func TestCreate(t *testing.T) {
	mockModel := FileUserModel{Filepath: testFilePath}

	currUsers, _ := mockModel.GetAll()
	oldLength := len(currUsers)
	var newID float64 = math.MinInt16
	for _, v := range currUsers {
		newID = math.Max(newID, float64(v.ID))
	}
	newID++

	testUser := model.User{FirstName: "testxyz", LastName: "ln", Email: "fake@email.org", Organization: "abc123"}

	mockModel.Create(&testUser)

	if testUser.ID != int(newID) {
		t.Errorf("expected ID %v got ID %v", newID, testUser.ID)
	}

	currUsers, _ = mockModel.GetAll()

	if len(currUsers) <= oldLength {
		t.Errorf("new user list should be longer than old user list")
	}

	newUser := getUserWithID(int(newID), currUsers)

	if newUser != testUser {
		t.Errorf("got %v want %v", newUser, testUser)
	}
}

func TestEdit(t *testing.T) {
	mockModel := FileUserModel{Filepath: testFilePath}

	currUsers, _ := mockModel.GetAll()

	editUser := currUsers[0]
	editUser.Email = "xyz@123"
	editUser.LastName = "zzzzz"

	mockModel.Edit(editUser, editUser.ID)

	currUsers, _ = mockModel.GetAll()

	storedEdits := getUserWithID(editUser.ID, currUsers)

	if editUser != storedEdits {
		t.Errorf("edits failed, got %v want %v", storedEdits, editUser)
	}
}

func TestDelete(t *testing.T) {
	mockModel := FileUserModel{Filepath: testFilePath}

	currUsers, _ := mockModel.GetAll()
	oldLength := len(currUsers)

	delID := currUsers[len(currUsers)-1].ID

	mockModel.Delete(delID)

	currUsers, _ = mockModel.GetAll()

	if len(currUsers) >= oldLength {
		t.Errorf("new user list should be shorter than old user list")
	}

	emptyUser := getUserWithID(delID, currUsers)
	testVal := model.User{}

	if emptyUser != testVal {
		t.Errorf("expected no user to be found but found %v", emptyUser)
	}
}

func MissingDelete(t *testing.T) {

}

func MissingEdit(t *testing.T) {

}

func IncompleteCreate(t *testing.T) {

}

func getUserWithID(ID int, uList []model.User) model.User {
	for _, v := range uList {
		if v.ID == ID {
			return v
		}
	}
	return model.User{}
}
