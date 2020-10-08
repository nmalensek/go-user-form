package validation

import (
	"testing"

	"github.com/nmalensek/go-user-form/model"
)

func TestValidEmails(t *testing.T) {
	emails := []string{"test@email.com", "a@bc.de", "hello-world@example.org", "first.last@co.co"}

	for _, v := range emails {
		if !EmailPattern.MatchString(v) {
			t.Errorf("Valid email %v did not match the regex.", v)
		}
	}
}

func TestInvalidEmails(t *testing.T) {
	emails := []string{"", "a@b.", "stringwithoutat", "test,,,@@mail.net"}

	for _, v := range emails {
		if EmailPattern.MatchString(v) {
			t.Errorf("Invalid email %v incorrectly matched the regex.", v)
		}
	}
}

func TestValidateInput(t *testing.T) {
	goodUser := model.User{FirstName: "test1", LastName: "Last",
		Email: "test@email.com", Organization: "Sales"}

	errors := ValidateCompleteInput(goodUser)

	if len(errors) != 0 {
		for _, e := range errors {
			t.Errorf("Error occurred when none was expected: %v", e.Message)
		}
	}
}

func TestBadEmail(t *testing.T) {
	badEmail := model.User{Email: "aaaaa"}

	errors := ValidateCompleteInput(badEmail)

	if len(errors) == 0 {
		t.Errorf("Input bad email address %v but was not caught", badEmail.Email)
	}

	want := IncorrectFormatMessage("Email")
	var got string
	for _, e := range errors {
		if e.PropName == "Email" {
			got = e.Message
			break
		}
	}
	if got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}

func TestIncompleteEntry(t *testing.T) {
	missingInfo := model.User{FirstName: "test"}

	errors := ValidateCompleteInput(missingInfo)

	if len(errors) < 3 {
		t.Errorf("Information missing from input was not caught.")
		for _, e := range errors {
			t.Errorf("Captured error: %v", e.PropName)
		}
	}

	names := make(map[string]struct{})
	names["LastName"] = struct{}{}
	names["Email"] = struct{}{}
	names["Organization"] = struct{}{}

	for _, e := range errors {
		if _, ok := names[e.PropName]; !ok {
			t.Errorf("%v is missing but no error was created for it.", e.PropName)
		}
	}
}

func TestErrorType(t *testing.T) {
	errs := UserErrors{Message: "Test"}

	result := checkType(errs)

	if result != "UserErrors" {
		t.Errorf("Expected type UserErrors but got %v", result)
	}
}

func checkType(e error) string {
	switch e.(type) {
	case UserErrors:
		return "UserErrors"
	default:
		return string(e.Error())
	}
}
