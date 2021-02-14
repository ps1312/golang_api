package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type FoodsStoreStub struct {
	foods []Food
}

func (f *FoodsStoreStub) GetFoods() ([]Food, error) {
	return f.foods, nil
}

func (f *FoodsStoreStub) PostFood(food Food) (Food, error) {
	return food, nil
}

type FailureStubStore struct{}

func (f *FailureStubStore) GetFoods() ([]Food, error) {
	return nil, errors.New(ErrInternalServer)
}

func (f *FailureStubStore) PostFood(food Food) (Food, error) {
	return Food{}, errors.New(ErrInternalServer)
}

type FoodsStoreSpy struct {
	calls          int
	postFoodParams Food
}

func (f *FoodsStoreSpy) GetFoods() ([]Food, error) {
	return nil, nil
}

func (f *FoodsStoreSpy) PostFood(food Food) (Food, error) {
	f.calls++
	f.postFoodParams = food

	return Food{}, nil
}

func TestGetFoods(t *testing.T) {
	server := &FoodsServer{}

	makeGetFoodsRequest := func() *http.Request {
		request, _ := http.NewRequest(http.MethodGet, "/foods", nil)
		return request
	}

	t.Run("returns multiple Foods in store", func(t *testing.T) {
		wantedFoods := []Food{{"food name 1", 300}, {"food name 2", 400}}
		server.store = &FoodsStoreStub{wantedFoods}
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

func TestPostFood(t *testing.T) {
	server := &FoodsServer{}

	makePostFoodRequest := func() *http.Request {
		var jsonStr = []byte(`{"name":"test","calories":111}`)
		request, _ := http.NewRequest(http.MethodPost, "/foods", bytes.NewBuffer(jsonStr))
		return request
	}

	t.Run("Delivers error on failure", func(t *testing.T) {
		server.store = &FailureStubStore{}
		response := httptest.NewRecorder()

		server.ServeHTTP(response, makePostFoodRequest())

		got := response.Body.String()
		want := ErrInternalServer

		assertStatus(t, response.Code, http.StatusInternalServerError)
		assertError(t, got, want)
	})

	t.Run("Delivers correct params to storage", func(t *testing.T) {
		wantedFood := Food{"test", 111}
		spy := &FoodsStoreSpy{}
		server.store = spy
		response := httptest.NewRecorder()

		server.ServeHTTP(response, makePostFoodRequest())

		if spy.calls != 1 {
			t.Errorf("got %d, want 1 call", spy.calls)
		}

		if spy.postFoodParams != wantedFood {
			t.Errorf("got %v, want %v", spy.postFoodParams, wantedFood)
		}
	})

	t.Run("Delivers missing params error on no Calories provided", func(t *testing.T) {
		server.store = &FoodsStoreStub{}
		body := strings.NewReader(`
		{
			"name": "food name 1"
		}
		`)
		request, _ := http.NewRequest(http.MethodPost, "/foods", body)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := ErrMissingParam

		assertStatus(t, response.Code, http.StatusUnprocessableEntity)
		assertError(t, got, want)
	})

	t.Run("Delivers created food and created status code", func(t *testing.T) {
		want := Food{"test", 111}
		server.store = &FoodsStoreStub{}
		response := httptest.NewRecorder()

		server.ServeHTTP(response, makePostFoodRequest())

		var got Food
		json.NewDecoder(response.Body).Decode(&got)
		assertStatus(t, response.Code, http.StatusCreated)

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
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
