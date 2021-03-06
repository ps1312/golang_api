package food

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
		server.Store = &FoodsStoreStub{wantedFoods}
		response := httptest.NewRecorder()

		server.ServeHTTP(response, makeGetFoodsRequest())

		assertStatus(t, response.Code, http.StatusOK)
		assertJSONBody(t, response.Body, wantedFoods)
	})

	t.Run("delivers 500 status code on storage error", func(t *testing.T) {
		server.Store = &FailureStubStore{}
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

	makePostFoodRequest := func(body string) *http.Request {
		request, _ := http.NewRequest(http.MethodPost, "/foods", strings.NewReader(body))
		return request
	}

	t.Run("Delivers error on failure", func(t *testing.T) {
		server.Store = &FailureStubStore{}
		response := httptest.NewRecorder()

		server.ServeHTTP(response, makePostFoodRequest(`{"name": "any-name","calories":123}`))

		got := response.Body.String()
		want := ErrInternalServer

		assertStatus(t, response.Code, http.StatusInternalServerError)
		assertError(t, got, want)
	})

	t.Run("Delivers correct params to storage", func(t *testing.T) {
		spy := &FoodsStoreSpy{}
		server.Store = spy
		want := Food{"test", 111}
		body := fmt.Sprintf(`{"name": %q,"calories":%d}`, want.Name, want.Calories)

		server.ServeHTTP(httptest.NewRecorder(), makePostFoodRequest(body))

		assertCallsCount(t, spy.calls, 1)

		if spy.postFoodParams != want {
			t.Errorf("got %v, want %v", spy.postFoodParams, want)
		}
	})

	t.Run("Delivers missing params error on no Calories provided", func(t *testing.T) {
		server.Store = &FoodsStoreStub{}
		body := `{"name": "food name 1"}`
		response := httptest.NewRecorder()

		server.ServeHTTP(response, makePostFoodRequest(body))

		assertStatus(t, response.Code, http.StatusUnprocessableEntity)
		assertMissingParam(t, response.Body.String(), "Calories")
	})

	t.Run("Delivers missing params error on no Name provided", func(t *testing.T) {
		server.Store = &FoodsStoreStub{}
		body := `{"calories": 1234}`
		response := httptest.NewRecorder()

		server.ServeHTTP(response, makePostFoodRequest(body))

		assertStatus(t, response.Code, http.StatusUnprocessableEntity)
		assertMissingParam(t, response.Body.String(), "Name")
	})

	t.Run("Delivers missing params error on invalid body", func(t *testing.T) {
		server.Store = &FoodsStoreStub{}
		body := `{"any_key": false}`
		response := httptest.NewRecorder()

		server.ServeHTTP(response, makePostFoodRequest(body))

		assertStatus(t, response.Code, http.StatusUnprocessableEntity)
		assertMissingParam(t, response.Body.String(), "Name")
	})

	t.Run("Does not call store on missing params error", func(t *testing.T) {
		spy := &FoodsStoreSpy{}
		server.Store = spy
		body := fmt.Sprintf(`{"name": %q}`, "any-name")

		server.ServeHTTP(httptest.NewRecorder(), makePostFoodRequest(body))

		assertCallsCount(t, spy.calls, 0)
	})

	t.Run("Delivers created food and created status code", func(t *testing.T) {
		server.Store = &FoodsStoreStub{}
		response := httptest.NewRecorder()

		want := Food{"test", 111}
		body := fmt.Sprintf(`{"name":%q,"calories":%d}`, want.Name, want.Calories)

		server.ServeHTTP(response, makePostFoodRequest(body))

		var got Food
		json.NewDecoder(response.Body).Decode(&got)
		assertStatus(t, response.Code, http.StatusCreated)

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

}

func assertCallsCount(t *testing.T, got int, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
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

func assertMissingParam(t *testing.T, got string, want string) {
	err := ErrMissingParam(want)
	if got != err.Error() {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func assertError(t *testing.T, got string, want string) {
	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}
