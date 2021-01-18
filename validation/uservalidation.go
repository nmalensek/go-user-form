package validation

import (
	"fmt"
	"regexp"

	"github.com/nmalensek/go-user-form/model"
)

//EmailPattern is a regular expression that should capture valid email addresses, source:
//https://html.spec.whatwg.org/multipage/input.html#valid-e-mail-address
var EmailPattern = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)

//MissingAllProps is an error that occurs if no properties are submitted to a PUT request.
const MissingAllProps = "At least one property must be filled out to complete an edit."

//ValidateCompleteInput compares a given user's properties against requirements and returns
//errors corresponding to the properties not meeting the requirements.
func ValidateCompleteInput(subj model.User) []UserError {
	errs := make([]UserError, 0)
	stringProps := getRequiredProps(&subj)

	errs = append(errs, checkIfRequired(stringProps)...)

	if len(subj.Email) > 0 {
		errs = append(errs, appendErrorIfNoMatch([]ValItem{
			{name: "Email", val: subj.Email, pattern: EmailPattern},
		})...)
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

//ValidatePartialInput checks the user object's property values. If at least one is filled out,
//then the input is ok. If the email property is included, it is additionally checked for formatting.
func ValidatePartialInput(subj model.User) []UserError {
	errs := make([]UserError, 0)
	reqProps := getRequiredProps(&subj)

	ok := userHasValue(reqProps)
	if !ok {
		errs = append(errs, UserError{PropName: "", PropValue: "", Message: MissingAllProps})
		return errs
	}

	if len(subj.Email) > 0 {
		appendErrorIfNoMatch([]ValItem{
			{name: "Email", val: subj.Email, pattern: EmailPattern}})
	}

	return nil
}

//ValItem is a property name and its value.
type ValItem struct {
	name         string
	val          string
	friendlyName string
	pattern      *regexp.Regexp
}

func (v *ValItem) getFriendlyName() string {
	if v.friendlyName == "" {
		return v.name
	}
	return v.friendlyName
}

//checkIfRequired appends an error if an item is missing and is required.
func checkIfRequired(items []ValItem) []UserError {
	e := make([]UserError, 0)
	for _, p := range items {
		if p.val == "" {
			uErr := UserError{PropName: p.name, PropValue: p.val, Message: RequiredMessage(p.getFriendlyName())}
			e = append(e, uErr)
		}
	}
	return e
}

//appendErrorIfNoMatch appends an error if an item does not match the specified regex.
func appendErrorIfNoMatch(items []ValItem) []UserError {
	e := make([]UserError, 0)
	for _, p := range items {
		if !p.pattern.MatchString(p.val) {
			newErr := UserError{PropName: p.name, PropValue: p.val, Message: IncorrectFormatMessage(p.getFriendlyName())}
			e = append(e, newErr)
		}
	}
	return e
}

func userHasValue(items []ValItem) bool {
	for _, p := range items {
		if p.val != "" {
			return true
		}
	}
	return false
}

func getRequiredProps(u *model.User) []ValItem {
	return []ValItem{
		{name: "FirstName", friendlyName: "First Name", val: u.FirstName},
		{name: "LastName", friendlyName: "Last Name", val: u.LastName},
		{name: "Email", val: u.Email},
		{name: "Organization", val: u.Organization}}
}

//IncorrectFormatMessage returns the standard error message for a property that's formatted incorrectly.
func IncorrectFormatMessage(prop string) string {
	return fmt.Sprintf("%v is not in the correct format.", prop)
}

//RequiredMessage returns the standard error message for a property that is required but is missing.
func RequiredMessage(prop string) string {
	return fmt.Sprintf("%v is required.", prop)
}
