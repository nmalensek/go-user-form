package validation

import (
	"fmt"
	"regexp"

	"github.com/nmalensek/go-user-form/model"
)

const emailPattern = `/^[a-zA-Z0-9.!#$%&'*+/=?^_{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/`

//ValidateInput compares a given user's properties against requirements and returns
//errors corresponding to the properties not meeting the requirements.
func ValidateInput(subj model.User) []UserError {
	errs := make([]UserError, 0)
	stringProps := []ValItem{
		ValItem{name: "firstName", val: subj.FirstName},
		ValItem{name: "lastName", val: subj.LastName},
		ValItem{name: "email", val: subj.Email},
		ValItem{name: "organization", val: subj.Organization}}

	appendErrorIfReq(errs, stringProps)

	appendErrorIfNoMatch(errs, []ValItem{
		ValItem{name: "email", val: subj.Email, pattern: emailPattern}})

	return errs
}

//UserError contains details about validation errors for a user object.
type UserError struct {
	PropName  string `json:"name"`
	PropValue string `json:"value"`
	Message   string `json:"msg"`
}

//ValItem is a property name and its value.
type ValItem struct {
	name    string
	val     string
	pattern string
}

func appendErrorIfReq(e []UserError, item []ValItem) {
	for _, p := range item {
		if p.val == "" {
			uErr := UserError{PropName: p.name, PropValue: p.val, Message: fmt.Sprintf("%v is required.", p.name)}
			e = append(e, uErr)
		}
	}
}

func appendErrorIfNoMatch(e []UserError, item []ValItem) {
	for _, p := range item {
		if match, err := regexp.MatchString(p.pattern, p.val); match == false || err != nil {
			newErr := UserError{PropName: p.name, PropValue: p.val, Message: fmt.Sprintf("%v is not in the correct format.", p.name)}
			e = append(e, newErr)
		}
	}
}
