package user

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type EncrypterSpy struct {
	calls           int
	encryptParam    string
	defaultPassword string
}

func (e *EncrypterSpy) encrypt(password string) string {
	e.calls++
	e.encryptParam = password
	return e.defaultPassword
}

func (e *EncrypterSpy) respondWith(password string) {
	e.defaultPassword = password
}

type UserStoreSpy struct {
	calls          int
	saveUserParams DatabaseModel
}

func (e *UserStoreSpy) save(user DatabaseModel) {
	e.calls++
	e.saveUserParams = user
}

func TestRegister(t *testing.T) {

	makeSUT := func() (UsersServer, *EncrypterSpy, *UserStoreSpy) {
		sut := UsersServer{}
		encrypter := &EncrypterSpy{}
		store := &UserStoreSpy{}
		sut.Encrypter = encrypter
		sut.Store = store

		return sut, encrypter, store
	}

	t.Run("Delivers 422 status code and missing param error correctly", func(t *testing.T) {
		sut, _, _ := makeSUT()

		want := ErrMissingParam("Name, Email, Password, PasswordConfirm")
		assertMissingParams(t, sut, nil, want.Error())

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
			assertMissingParams(t, sut, strings.NewReader(testCase.body), testCase.want.Error())
		}
	})

	t.Run("Delivers 422 status code and ErrPasswordsDontMatch error on not equal passwords", func(t *testing.T) {
		sut, _, _ := makeSUT()
		body := `{"name":"any-name", "email": "email@mail.com", "password": "password123", "passwordConfirm": "diffPassword"}`
		request, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
		response := httptest.NewRecorder()
		sut.ServeHTTP(response, request)

		got := response.Body.String()
		want := ErrPasswordsDontMatch

		assertStatusCode(t, response.Code, http.StatusUnprocessableEntity)
		assertError(t, got, want)
	})

	t.Run("Calls encrypter with correct password", func(t *testing.T) {
		sut, encrypter, _ := makeSUT()
		body := `{"name":"any-name", "email": "email@mail.com", "password": "password123", "passwordConfirm": "password123"}`
		request, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
		response := httptest.NewRecorder()
		sut.ServeHTTP(response, request)

		assertCalls(t, encrypter.calls, 1)
		assertString(t, encrypter.encryptParam, "password123")
	})

	t.Run("Delivers internal server error on encryptor failure", func(t *testing.T) {
		sut, _, _ := makeSUT()
		body := `{"name":"any-name", "email": "email@mail.com", "password": "password123", "passwordConfirm": "password123"}`
		request, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
		response := httptest.NewRecorder()
		sut.ServeHTTP(response, request)

		assertStatusCode(t, response.Code, http.StatusInternalServerError)
		assertError(t, response.Body.String(), ErrInternalServer)
	})

	t.Run("Calls store with correct user and encrypted password", func(t *testing.T) {
		sut, encrypter, store := makeSUT()
		body := `{"name":"any-name", "email": "email@mail.com", "password": "password123", "passwordConfirm": "password123"}`
		request, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
		response := httptest.NewRecorder()

		const wantedEncryptedPassword = "hashed_password"
		encrypter.respondWith(wantedEncryptedPassword)

		sut.ServeHTTP(response, request)

		got := store.saveUserParams
		want := DatabaseModel{Name: "any-name", Email: "email@mail.com", Password: wantedEncryptedPassword}

		assertCalls(t, store.calls, 1)
		assertString(t, got.Name, want.Name)
		assertString(t, got.Email, want.Email)
		assertString(t, got.Password, wantedEncryptedPassword)
	})
}

func assertMissingParams(t *testing.T, sut UsersServer, body io.Reader, want string) {
	t.Helper()

	request, _ := http.NewRequest(http.MethodPost, "/register", body)
	response := httptest.NewRecorder()

	sut.ServeHTTP(response, request)

	got := response.Body.String()

	assertStatusCode(t, response.Code, http.StatusUnprocessableEntity)
	assertError(t, got, want)
}

func assertCalls(t *testing.T, got int, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func assertString(t *testing.T, got string, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
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
