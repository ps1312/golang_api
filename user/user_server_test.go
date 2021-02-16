package user

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type EncrypterSpy struct {
	calls        int
	encryptParam string
}

func (e *EncrypterSpy) encrypt(password string) string {
	e.calls++
	e.encryptParam = password
	return ""
}

func TestRegister(t *testing.T) {
	server := UsersServer{}

	t.Run("Delivers 422 status code and missing param error correctly", func(t *testing.T) {
		spy := &EncrypterSpy{}
		server.Encrypter = spy

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
		spy := &EncrypterSpy{}
		server.Encrypter = spy
		body := `{"name":"any-name", "email": "email@mail.com", "password": "password123", "passwordConfirm": "diffPassword"}`
		request, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := ErrPasswordsDontMatch

		assertStatusCode(t, response.Code, http.StatusUnprocessableEntity)
		assertError(t, got, want)
	})

	t.Run("Calls encrypter with correct password", func(t *testing.T) {
		spy := &EncrypterSpy{}
		server.Encrypter = spy
		body := `{"name":"any-name", "email": "email@mail.com", "password": "password123", "passwordConfirm": "password123"}`
		request, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		if spy.calls != 1 {
			t.Errorf("got %d, want 1", spy.calls)
		}

		if spy.encryptParam != "password123" {
			t.Errorf("got %q, want %q", spy.encryptParam, "password123")
		}
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
