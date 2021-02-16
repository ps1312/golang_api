package user

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrInternalServer error const
const ErrInternalServer = "Internal server error"

// ErrPasswordsDontMatch error const
const ErrPasswordsDontMatch = "Passwords don't match"

// ErrMissingParam error struct for displaying missing param error with specified param
type ErrMissingParam string

func (e *ErrMissingParam) Error() string {
	return fmt.Sprintf("Missing parameter(s): %q", *e)
}

// DatabaseModel user struct
type DatabaseModel struct {
	Name     string
	Email    string
	Password string
}

// RegisterModel model struct
type RegisterModel struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
}

// Encrypter password interface
type Encrypter interface {
	encrypt(password string) string
}

// Store user store interface
type Store interface {
	save(user DatabaseModel)
}

// UsersServer struct
type UsersServer struct {
	Encrypter Encrypter
	Store     Store
}

func (u *UsersServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		err := ErrMissingParam("Name, Email, Password, PasswordConfirm")
		respondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	var user RegisterModel
	json.NewDecoder(req.Body).Decode(&user)

	missingParams := ErrMissingParam(checkMissingParams(user))
	if missingParams != "" {
		respondWithError(w, http.StatusUnprocessableEntity, missingParams.Error())
		return
	}

	if user.Password != user.PasswordConfirm {
		respondWithError(w, http.StatusUnprocessableEntity, ErrPasswordsDontMatch)
		return
	}

	hashed := u.Encrypter.encrypt(user.Password)

	respondWithError(w, http.StatusInternalServerError, ErrInternalServer)
	u.Store.save(DatabaseModel{Name: user.Name, Email: user.Email, Password: hashed})
}

func respondWithError(w http.ResponseWriter, status int, err string) {
	w.WriteHeader(status)
	fmt.Fprint(w, err)
}

func checkMissingParams(user RegisterModel) (missingParams string) {
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
