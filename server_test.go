package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type FoodsStoreSpy struct {
	foods []Food
}

func (f *FoodsStoreSpy) GetFoods() []Food {
	return f.foods
}

func TestGetFoods(t *testing.T) {
	server := &FoodsServer{}

	makeGetFoodsRequest := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/foods", nil)
		return request
	}
	makeResponseRecorder := func() *httptest.ResponseRecorder {
		return httptest.NewRecorder()
	}

	t.Run("returns empty list on empty store", func(t *testing.T) {
		server.store = &FoodsStoreSpy{[]Food{}}
		response := makeResponseRecorder()

		server.ServeHTTP(response, makeGetFoodsRequest())

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "[]")
	})

	t.Run("returns array of one item on single Food in store", func(t *testing.T) {
		server.store = &FoodsStoreSpy{[]Food{{name: "food 1", calories: 300}}}
		response := makeResponseRecorder()

		server.ServeHTTP(response, makeGetFoodsRequest())

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "[{food 1 300}]")
	})

}

func assertStatus(t *testing.T, got int, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertResponseBody(t *testing.T, got string, want string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}
