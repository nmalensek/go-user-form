package users

import (
	"net/http"

	"github.com/nmalensek/go-user-form/config"
)

//ProcessRequestByType checks which HTTP verb the request has and processes it accordingly.
func ProcessRequestByType(w http.ResponseWriter, r *http.Request, e *config.Env) {

}
