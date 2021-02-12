package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type EmptyFoodsStoreSpy struct {
}

func (f *EmptyFoodsStoreSpy) GetFoods() string {
	return "[]"
}

func TestGetFoods(t *testing.T) {
	t.Run("returns empty list on empty store", func(t *testing.T) {
		server := &FoodsServer{store: &EmptyFoodsStoreSpy{}}
		request, _ := http.NewRequest(http.MethodGet, "/foods", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "[]"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
