package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type FoodsStoreSpy struct {
	foods []Food
}

func (f *FoodsStoreSpy) GetFoods() ([]Food, error) {
	return f.foods, nil
}

type FailureStubStore struct{}

func (f *FailureStubStore) GetFoods() ([]Food, error) {
	return nil, errors.New(ErrInternalServer)
}

func TestGetFoods(t *testing.T) {
	server := &FoodsServer{}

	makeGetFoodsRequest := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/foods", nil)
		return request
	}

	t.Run("returns multiple Foods in store", func(t *testing.T) {
		wantedFoods := []Food{{"food name 1", 300}, {"food name 2", 400}}
		server.store = &FoodsStoreSpy{wantedFoods}
		response := httptest.NewRecorder()

		server.ServeHTTP(response, makeGetFoodsRequest())

		assertStatus(t, response.Code, http.StatusOK)
		assertJSONBody(t, response.Body, wantedFoods)
	})

	t.Run("delivers 500 status code on storage error", func(t *testing.T) {
		server.store = &FailureStubStore{}
		response := httptest.NewRecorder()

		server.ServeHTTP(response, makeGetFoodsRequest())

		got := response.Body.String()
		want := ErrInternalServer

		assertStatus(t, response.Code, http.StatusInternalServerError)
		assertError(t, got, want)
	})
}

func assertStatus(t *testing.T, got int, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
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

func assertError(t *testing.T, got string, want string) {
	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}
