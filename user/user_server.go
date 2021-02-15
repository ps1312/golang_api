package user

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrPasswordsDontMatch error const
const ErrPasswordsDontMatch = "Passwords don't match."

// ErrMissingParam error struct for displaying missing param error with specified param
type ErrMissingParam string

func (e *ErrMissingParam) Error() string {
	return fmt.Sprintf("Missing parameter(s): %q", *e)
}

// User model struct
type User struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
}

// UsersServer struct
type UsersServer struct{}

func (u *UsersServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		err := ErrMissingParam("Name, Email, Password, PasswordConfirm")
		fmt.Fprint(w, err.Error())
		return
	}

	var user User
	json.NewDecoder(req.Body).Decode(&user)

	w.WriteHeader(http.StatusUnprocessableEntity)
	missing := ErrMissingParam(checkMissingParams(user))

	if missing != "" {
		fmt.Fprint(w, missing.Error())
		return
	}

	fmt.Fprint(w, ErrPasswordsDontMatch)
}

func checkMissingParams(user User) (missingParams string) {
	if user.Name == "" {
		missingParams += "Name, "
	}

	if user.Email == "" {
		missingParams += "Email, "
	}

	if user.Password == "" {
		missingParams += "Password, "
	}

	if user.PasswordConfirm == "" {
		missingParams += "PasswordConfirm, "
	}

	if missingParams != "" {
		missingParams = missingParams[:len(missingParams)-2]
	}
	return
}
