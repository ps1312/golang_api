package user

import (
	"api/encryption"
	"api/helpers"
	"api/signer"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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
	getAll() ([]DatabaseModel, error)
}

// Server struct
type Server struct {
	Encrypter encryption.Encrypter
	Store     Store
	Signer    signer.Signer
}

func (u *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		handleGetUsers(w, u.Store)
	} else {
		handlePostUser(w, req, u.Store, u.Encrypter, u.Signer)
	}
}

func handlePostUser(w http.ResponseWriter, req *http.Request, store Store, encryptor encryption.Encrypter, signer signer.Signer) {
	if req.Body == nil {
		err := ErrMissingParam("Name, Email, Password, PasswordConfirm")
		helpers.RespondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	var user RegisterModel
	json.NewDecoder(req.Body).Decode(&user)

	requiredFields := []string{"Name", "Email", "Password", "PasswordConfirm"}
	missingParams := ErrMissingParam(checkMissingParams(user, requiredFields))

	if missingParams != "" {
		helpers.RespondWithError(w, http.StatusUnprocessableEntity, missingParams.Error())
		return
	}

	if user.Password != user.PasswordConfirm {
		helpers.RespondWithError(w, http.StatusUnprocessableEntity, ErrPasswordsDontMatch)
		return
	}

	hashed, hashErr := encryptor.Encrypt(user.Password, 10)

	if hashErr != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, ErrInternalServer)
		return
	}

	dbUser := DatabaseModel{Name: user.Name, Email: user.Email, password: hashed}
	storeErr := store.save(dbUser)

	if storeErr != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, ErrInternalServer)
		return
	}

	exp := time.Now().Add(time.Minute * 15).Unix()
	token, signerErr := signer.Sign(user.Name, user.Email, exp)

	if signerErr != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, ErrInternalServer)
		return
	}

	signedUser := struct {
		User  DatabaseModel
		Token string
	}{dbUser, token}

	helpers.RespondWithSuccess(w, http.StatusCreated, signedUser)
}

func handleGetUsers(w http.ResponseWriter, store Store) {
	users, err := store.getAll()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, ErrInternalServer)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users)
	}

}

func checkMissingParams(user RegisterModel, params []string) (missingParams string) {
	for _, param := range params {
		if helpers.GetFieldFromStruct(user, param) == "" {
			missingParams += param + ", "
		}
	}
	if missingParams != "" {
		missingParams = missingParams[:len(missingParams)-2]
	}
	return
}
