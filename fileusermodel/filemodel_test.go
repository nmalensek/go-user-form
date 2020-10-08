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

	editData := model.User{Email: "xyz@123", LastName: "zzzzz", Organization: "..."}
	originalUser := currUsers[0]

	originalUser.Email = "xyz@123"
	originalUser.LastName = "zzzzz"
	originalUser.Organization = "..."

	mockModel.Edit(editData, originalUser.ID)

	currUsers, _ = mockModel.GetAll()

	storedEdits := getUserWithID(originalUser.ID, currUsers)

	if originalUser != storedEdits {
		t.Errorf("edit failed, got %v want %v", storedEdits, originalUser)
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

func TestMissingDelete(t *testing.T) {
	mockModel := FileUserModel{Filepath: testFilePath}
	err := mockModel.Delete(math.MaxInt64)

	if err == nil {
		t.Errorf("expected error, got none.")
		return
	}

	if err.Error() != model.CouldNotFind {
		t.Errorf("error mismatch; got %v want %v", err.Error(), model.CouldNotFind)
	}
}

func TestMissingEdit(t *testing.T) {
	//covered by MissingDelete
}

func TestIncompleteCreate(t *testing.T) {
	mockModel := FileUserModel{Filepath: testFilePath}

	incompleteUser := model.User{FirstName: "test"}

	err := mockModel.Create(&incompleteUser)

	if err == nil {
		t.Errorf("expected error, got none.")
		return
	}

	if err.Error() != model.CreateErrorIncomplete {
		t.Errorf("error mismatch; got %v want %v", err.Error(), model.CouldNotFind)
	}
}

func getUserWithID(ID int, uList []model.User) model.User {
	for _, v := range uList {
		if v.ID == ID {
			return v
		}
	}
	return model.User{}
}
