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

// Encrypter password interface
type Encrypter interface {
	encrypt(password string) string
}

// UsersServer struct
type UsersServer struct {
	Encrypter Encrypter
}

func (u *UsersServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		err := ErrMissingParam("Name, Email, Password, PasswordConfirm")
		respondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	var user User
	json.NewDecoder(req.Body).Decode(&user)

	missingParams := ErrMissingParam(checkMissingParams(user))
	if missingParams != "" {
		respondWithError(w, http.StatusUnprocessableEntity, missingParams.Error())
		return
	}

	u.Encrypter.encrypt(user.Password)
	fmt.Fprint(w, ErrPasswordsDontMatch)
}

func respondWithError(w http.ResponseWriter, status int, err string) {
	w.WriteHeader(status)
	fmt.Fprint(w, err)
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
