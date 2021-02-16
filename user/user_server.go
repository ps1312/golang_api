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
	password string
}

// RegisterModel model struct
type RegisterModel struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
}

// Store user store interface
type Store interface {
	save(user DatabaseModel) error
}

// Server struct
type Server struct {
	Encrypter Encrypter
	Store     Store
}

func (u *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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

	hashed, hashErr := u.Encrypter.encrypt(user.Password, 10)

	if hashErr != nil {
		respondWithError(w, http.StatusInternalServerError, ErrInternalServer)
		return
	}

	dbUser := DatabaseModel{Name: user.Name, Email: user.Email, password: hashed}
	storeErr := u.Store.save(dbUser)

	if storeErr != nil {
		respondWithError(w, http.StatusInternalServerError, ErrInternalServer)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dbUser)
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
