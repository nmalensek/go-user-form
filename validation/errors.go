package validation

import "encoding/json"

//UserError contains details about validation errors for a user object.
type UserError struct {
	PropName  string `json:"name"`
	PropValue string `json:"value"`
	Message   string `json:"msg"`
}

//Error returns a JSONified version of a UserError.
func (u *UserError) Error() string {
	res, err := json.Marshal(u)
	if err != nil {
		return err.Error()
	}
	return string(res)
}

//UserErrors contains details about what error occurred and a slice of specific errors.
type UserErrors struct {
	Message   string
	ErrorList []UserError `json:"errors"`
}

func (u UserErrors) Error() string {
	return u.Message
}
