package fileusermodel

import (
	"io/ioutil"
	"os"
	"testing"
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
