package user

import (
	"encoding/json"
	"errors"
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
	defaultError    error
}

func (e *EncrypterSpy) Encrypt(password string, cost int) (string, error) {
	e.calls++
	e.encryptParam = password
	return e.defaultPassword, e.defaultError
}

func (e *EncrypterSpy) respondWith(password string) {
	e.defaultPassword = password
}

func (e *EncrypterSpy) respondWithError(err error) {
	e.defaultError = err
}

type UserStoreSpy struct {
	calls          int
	saveUserParams DatabaseModel
	defaultError   error
	Users          []DatabaseModel
}

func (e *UserStoreSpy) save(user DatabaseModel) error {
	e.calls++
	e.saveUserParams = user
	return e.defaultError
}

func (e *UserStoreSpy) getAll() ([]DatabaseModel, error) {
	return e.Users, e.defaultError
}

func (e *UserStoreSpy) respondGetAllWith(users []DatabaseModel) {
	e.Users = users
}

func (e *UserStoreSpy) respondWithError(err error) {
	e.defaultError = err
}

type SignerSpy struct {
	defaultError error
	defaultToken string
}

func (s *SignerSpy) respondWith(token string) {
	s.defaultToken = token
}

func (s *SignerSpy) Sign() (string, error) {
	return s.defaultToken, s.defaultError
}

func TestRegister(t *testing.T) {
	t.Run("Delivers 422 status code and missing param error correctly", func(t *testing.T) {
		sut, _, _, _ := makeSUT(t)

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
		body := `{"name":"any-name", "email": "email@mail.com", "password": "password123", "passwordConfirm": "diffPassword"}`
		sut, _, _, _ := makeSUT(t)

		response := makeRequestForRegistration(t, sut, body)

		got := response.Body.String()
		want := ErrPasswordsDontMatch

		assertStatusCode(t, response.Code, http.StatusUnprocessableEntity)
		assertError(t, got, want)
	})

	t.Run("Calls encrypter with correct password", func(t *testing.T) {
		sut, encrypter, _, _ := makeSUT(t)
		makeRequestForRegistration(t, sut, makeValidBody())

		assertCalls(t, encrypter.calls, 1)
		assertString(t, encrypter.encryptParam, "password123")
	})

	t.Run("Delivers internal server error on encryptor failure", func(t *testing.T) {
		sut, encrypter, _, _ := makeSUT(t)
		encrypter.respondWithError(errors.New("any-error"))

		response := makeRequestForRegistration(t, sut, makeValidBody())

		assertStatusCode(t, response.Code, http.StatusInternalServerError)
		assertError(t, response.Body.String(), ErrInternalServer)
	})

	t.Run("Calls store with correct user and encrypted password", func(t *testing.T) {
		sut, encrypter, store, _ := makeSUT(t)
		const wantedEncryptedPassword = "hashed_password"
		encrypter.respondWith(wantedEncryptedPassword)

		makeRequestForRegistration(t, sut, makeValidBody())

		got := store.saveUserParams
		want := DatabaseModel{Name: "any-name", Email: "email@mail.com", password: wantedEncryptedPassword}

		assertCalls(t, store.calls, 1)
		assertString(t, got.Name, want.Name)
		assertString(t, got.Email, want.Email)
		assertString(t, got.password, wantedEncryptedPassword)
	})

	t.Run("Delivers 500 status code on store error", func(t *testing.T) {
		sut, encrypter, store, _ := makeSUT(t)
		store.respondWithError(errors.New("any-error"))
		encrypter.respondWith("hashed_password")

		response := makeRequestForRegistration(t, sut, makeValidBody())

		assertStatusCode(t, response.Code, http.StatusInternalServerError)
		assertError(t, response.Body.String(), ErrInternalServer)
	})

	t.Run("Delivers error on token creation failure", func(t *testing.T) {
		sut, _, _, signer := makeSUT(t)
		signer.defaultError = errors.New("any-error")

		response := makeRequestForRegistration(t, sut, makeValidBody())

		assertStatusCode(t, response.Code, http.StatusInternalServerError)
		assertString(t, response.Body.String(), ErrInternalServer)
	})

	t.Run("Delivers 201 status code and created user without password", func(t *testing.T) {
		sut, encrypter, _, signer := makeSUT(t)
		encrypter.respondWith("hashed_password")
		signer.respondWith("signed_token")

		response := makeRequestForRegistration(t, sut, makeValidBody())

		got := response.Body.String()
		want := `{"User":{"Name":"any-name","Email":"email@mail.com"},"Token":"signed_token"}` + "\n"

		assertStatusCode(t, response.Code, http.StatusCreated)
		assertString(t, got, want)
	})
}

func TestGetUsers(t *testing.T) {
	t.Run("Delivers error on storage failure", func(t *testing.T) {
		sut, _, store, _ := makeSUT(t)
		store.respondWithError(errors.New(ErrInternalServer))
		request, _ := http.NewRequest(http.MethodGet, "/users", nil)
		response := httptest.NewRecorder()

		sut.ServeHTTP(response, request)

		assertStatusCode(t, response.Code, http.StatusInternalServerError)
		assertError(t, response.Body.String(), ErrInternalServer)
	})

	t.Run("Delivers users in storage successfully", func(t *testing.T) {
		sut, _, store, _ := makeSUT(t)
		want := []DatabaseModel{{Name: "any-Name", Email: "mail@mail.com", password: "any-password"}}
		store.respondGetAllWith(want)
		store.Users = want

		request, _ := http.NewRequest(http.MethodGet, "/users", nil)
		response := httptest.NewRecorder()

		sut.ServeHTTP(response, request)

		var got []DatabaseModel
		json.NewDecoder(response.Body).Decode(&got)
		assertStatusCode(t, response.Code, http.StatusOK)
		assertString(t, got[0].Name, want[0].Name)
		assertString(t, got[0].Email, want[0].Email)
	})
}

func makeSUT(t *testing.T) (Server, *EncrypterSpy, *UserStoreSpy, *SignerSpy) {
	sut := Server{}
	encrypter := &EncrypterSpy{}
	store := &UserStoreSpy{}
	signer := &SignerSpy{}
	signer.defaultError = nil
	sut.Encrypter = encrypter
	sut.Store = store
	sut.Signer = signer

	return sut, encrypter, store, signer
}

func makeRequestForRegistration(t *testing.T, sut Server, body string) httptest.ResponseRecorder {
	request, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	response := httptest.NewRecorder()
	sut.ServeHTTP(response, request)

	return *response
}

func makeValidBody() string {
	return `{"name":"any-name", "email": "email@mail.com", "password": "password123", "passwordConfirm": "password123"}`

}

func assertMissingParams(t *testing.T, sut Server, body io.Reader, want string) {
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
