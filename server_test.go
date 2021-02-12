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
	t.Run("returns empty list on empty store", func(t *testing.T) {
		foods := []Food{}
		server := &FoodsServer{store: &FoodsStoreSpy{foods}}
		request, _ := http.NewRequest(http.MethodGet, "/foods", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "[]"

		assertStatus(t, response.Code, http.StatusOK)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("returns array of one item on single Food in store", func(t *testing.T) {
		foods := []Food{{name: "food 1", calories: 300}}
		server := FoodsServer{store: &FoodsStoreSpy{foods}}
		request, _ := http.NewRequest(http.MethodGet, "/foods", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "[{food 1 300}]"

		assertStatus(t, response.Code, http.StatusOK)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func assertStatus(t *testing.T, got int, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
