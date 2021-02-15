package user

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// ErrMissingParam error struct for displaying missing param error with specified param
type ErrMissingParam string

func (e *ErrMissingParam) Error() string {
	return fmt.Sprintf("Missing parameter(s): %q", *e)
}

type User struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
}

type UsersServer struct{}

func (u *UsersServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		err := ErrMissingParam("Name, Email, Password, PasswordConfirm")
		fmt.Fprint(w, err.Error())
		return
	}

	missingParams := ""
	var user User
	json.NewDecoder(req.Body).Decode(&user)

	if user.Password == "" {
		missingParams = "Name, Email, Password, PasswordConfirm"
	} else {
		missingParams = "Name"
	}

	w.WriteHeader(http.StatusUnprocessableEntity)
	err := ErrMissingParam(missingParams)
	fmt.Fprint(w, err.Error())
}

func TestRegister(t *testing.T) {
	server := UsersServer{}

	t.Run("Delivers 422 status code and missing param error on no body provided", func(t *testing.T) {
		want := ErrMissingParam("Name, Email, Password, PasswordConfirm")
		assertMissingParams(t, server, nil, want.Error())
	})

	t.Run("Delivers 422 status code and missing param error on no params provided", func(t *testing.T) {
		want := ErrMissingParam("Name, Email, Password, PasswordConfirm")
		assertMissingParams(t, server, strings.NewReader(""), want.Error())
	})

	t.Run("Delivers 422 status code and missing params error on no Name provided", func(t *testing.T) {
		body := `{"email": "email@mail.com", "password": "password123", "passwordConfirm": "password123"}`
		want := ErrMissingParam("Name")
		assertMissingParams(t, server, strings.NewReader(body), want.Error())
	})
}

func assertMissingParams(t *testing.T, server UsersServer, body io.Reader, want string) {
	t.Helper()

	request, _ := http.NewRequest(http.MethodPost, "/register", body)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	got := response.Body.String()

	assertStatusCode(t, response.Code, http.StatusUnprocessableEntity)
	assertError(t, got, want)
}

func assertError(t *testing.T, got string, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertStatusCode(t *testing.T, got int, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
