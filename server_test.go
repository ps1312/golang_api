package main

import (
	"bytes"
	"encoding/json"
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

	t.Run("returns empty list on empty store", func(t *testing.T) {
		wantedFoods := []Food{}
		server.store = &FoodsStoreSpy{wantedFoods}
		response := httptest.NewRecorder()

		server.ServeHTTP(response, makeGetFoodsRequest())

		assertStatus(t, response.Code, http.StatusOK)
		assertJSONBody(t, response.Body, wantedFoods)
	})

	t.Run("returns single Food in store", func(t *testing.T) {
		wantedFoods := []Food{{"food name 1", 300}}
		server.store = &FoodsStoreSpy{wantedFoods}
		response := httptest.NewRecorder()

		server.ServeHTTP(response, makeGetFoodsRequest())

		assertStatus(t, response.Code, http.StatusOK)
		assertJSONBody(t, response.Body, wantedFoods)
	})
}

func assertStatus(t *testing.T, got int, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertJSONBody(t *testing.T, body *bytes.Buffer, want []Food) {
	t.Helper()

	var got []Food
	err := json.NewDecoder(body).Decode(&got)

	if err != nil {
		t.Fatalf("Unable to decode: error %q", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}
