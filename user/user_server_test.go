package user

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegister(t *testing.T) {
	server := UsersServer{}

	t.Run("Delivers 422 status code and missing param error correctly", func(t *testing.T) {
		want := ErrMissingParam("Name, Email, Password, PasswordConfirm")
		assertMissingParams(t, server, nil, want.Error())

		testingTable := []struct {
			body string
			want ErrMissingParam
		}{
			{"", ErrMissingParam("Name, Email, Password, PasswordConfirm")},
			{`{"email": "email@mail.com", "password": "password123", "passwordConfirm": "password123"}`, ErrMissingParam("Name")},
			{`{"name":"any-name", "password": "password123", "passwordConfirm": "password123"}`, ErrMissingParam("Email")},
			{`{"name":"any-name", "email": "email@mail.com", "passwordConfirm": "password123"}`, ErrMissingParam("Password")},
			{`{"name":"any-name", "email": "email@mail.com", "password": "password123"}`, ErrMissingParam("PasswordConfirm")},
		}

		for _, testCase := range testingTable {
			assertMissingParams(t, server, strings.NewReader(testCase.body), testCase.want.Error())
		}
	})

	t.Run("Delivers 422 status code and ErrPasswordsDontMatch error on not equal passwords", func(t *testing.T) {
		body := `{"name":"any-name", "email": "email@mail.com", "password": "password123", "passwordConfirm": "diffPassword"}`
		request, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := ErrPasswordsDontMatch

		assertStatusCode(t, response.Code, http.StatusUnprocessableEntity)
		assertError(t, got, want)
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
