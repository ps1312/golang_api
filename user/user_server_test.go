package user

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ErrMissingParam error struct for displaying missing param error with specified param
type ErrMissingParam string

func (e *ErrMissingParam) Error() string {
	return fmt.Sprintf("Missing parameter: %q", *e)
}

type UsersServer struct{}

func (u *UsersServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	err := ErrMissingParam("Name, Email, Password, PasswordConfirm")
	fmt.Fprint(w, err.Error())
}

func TestRegister(t *testing.T) {
	server := UsersServer{}

	t.Run("Delivers 422 status code and missing param error on no params provided", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/register", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := ErrMissingParam("Name, Email, Password, PasswordConfirm")

		assertStatusCode(t, response.Code, http.StatusUnprocessableEntity)
		assertError(t, got, want.Error())
	})
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
