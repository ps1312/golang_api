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
	err := ErrMissingParam("Name")
	fmt.Fprint(w, err.Error())
}

func TestRegisterDelivers422StatusCodeAndMissingParamErrorOnNoParamsProvided(t *testing.T) {
	server := UsersServer{}
	request, _ := http.NewRequest(http.MethodPost, "/register", nil)
	response := httptest.NewRecorder()
	fmt.Println(ErrMissingParam("Name"))

	server.ServeHTTP(response, request)

	got := response.Body.String()
	want := ErrMissingParam("Name")

	code := response.Code
	if code != http.StatusUnprocessableEntity {
		t.Errorf("got %d, want 422", code)
	}

	if got != want.Error() {
		t.Errorf("got %v, want %v", got, want)
	}
}
