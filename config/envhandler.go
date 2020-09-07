package config

import "net/http"

//MakeHandler checks the requested path and returns 404 if not found. Otherwise, it calls the handler function passed in that requires an environment variable.
func MakeHandler(fn func(w http.ResponseWriter, r *http.Request, e *Env), env *Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, env)
	}
}
