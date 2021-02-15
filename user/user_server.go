package user

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
	err := ErrMissingParam(checkMissingParams(user))
	fmt.Fprint(w, err.Error())
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
	missingParams = missingParams[:len(missingParams)-2]
	return
}
