package validation

import (
	"fmt"
	"regexp"

	"github.com/nmalensek/go-user-form/model"
)

//EmailPattern is a regular expression that should capture valid email addresses, source:
//https://html.spec.whatwg.org/multipage/input.html#valid-e-mail-address
var EmailPattern = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)

//ValidateInput compares a given user's properties against requirements and returns
//errors corresponding to the properties not meeting the requirements.
func ValidateInput(subj model.User) []UserError {
	errs := make([]UserError, 0)
	stringProps := []ValItem{
		ValItem{name: "FirstName", val: subj.FirstName},
		ValItem{name: "LastName", val: subj.LastName},
		ValItem{name: "Email", val: subj.Email},
		ValItem{name: "Organization", val: subj.Organization}}

	appendErrorIfReq(&errs, stringProps)

	if len(subj.Email) > 0 {
		appendErrorIfNoMatch(&errs, []ValItem{
			ValItem{name: "Email", val: subj.Email, pattern: EmailPattern}})
	}

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
	pattern *regexp.Regexp
}

//appendErrorIfReq appends an error if an item is missing and is required.
func appendErrorIfReq(e *[]UserError, item []ValItem) {
	for _, p := range item {
		if p.val == "" {
			uErr := UserError{PropName: p.name, PropValue: p.val, Message: RequiredMessage(p.name)}
			*e = append(*e, uErr)
		}
	}
}

//appendErrorIfNoMatch appends an error if an item does not match the specified regex.
func appendErrorIfNoMatch(e *[]UserError, item []ValItem) {
	for _, p := range item {
		if !p.pattern.MatchString(p.val) {
			newErr := UserError{PropName: p.name, PropValue: p.val, Message: IncorrectFormatMessage(p.name)}
			*e = append(*e, newErr)
		}
	}
}

//IncorrectFormatMessage returns the standard error message for a property that's formatted incorrectly.
func IncorrectFormatMessage(prop string) string {
	return fmt.Sprintf("%v is not in the correct format.", prop)
}

//RequiredMessage returns the standard error message for a property that is required but is missing.
func RequiredMessage(prop string) string {
	return fmt.Sprintf("%v is required.", prop)
}
