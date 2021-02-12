package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type EmptyFoodsStoreSpy struct {
}

func (f *EmptyFoodsStoreSpy) GetFoods() []Food {
	return []Food{}
}

type SingleItemFoodsStore struct{}

func (f *SingleItemFoodsStore) GetFoods() []Food {
	foods := make([]Food, 0)
	foods = append(foods, Food{name: "food 1", calories: 300})
	return foods
}

func TestGetFoods(t *testing.T) {
	t.Run("returns empty list on empty store", func(t *testing.T) {
		server := &FoodsServer{store: &EmptyFoodsStoreSpy{}}
		request, _ := http.NewRequest(http.MethodGet, "/foods", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "[]"

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("returns array of one item on single Food in store", func(t *testing.T) {
		server := FoodsServer{store: &SingleItemFoodsStore{}}
		request, _ := http.NewRequest(http.MethodGet, "/foods", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "[{food 1 300}]"

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
