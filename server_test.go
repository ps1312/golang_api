package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func FoodServer(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "[]")
}

func TestGetFoods(t *testing.T) {
	t.Run("returns empty list", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/foods", nil)
		response := httptest.NewRecorder()

		FoodServer(response, request)

		got := response.Body.String()
		want := "[]"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
