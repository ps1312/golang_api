package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrInternalServer constant for error message
const ErrInternalServer = "Internal server error"

// Food struct type
type Food struct {
	Name     string `json:"name"`
	Calories int    `json:"calories"`
}

// FoodsStore interface for Food storage operations
type FoodsStore interface {
	GetFoods() ([]Food, error)
}

// FoodsServer struct to use FoodsStore
type FoodsServer struct {
	store FoodsStore
}

// FoodServer handles requests for foods
func (f *FoodsServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		foods, err := f.store.GetFoods()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, ErrInternalServer)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(foods)
		}
	}

	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, ErrInternalServer)
}

// InMemoryFoodsStore in memory store for testing
type InMemoryFoodsStore struct{}

// GetFoods returns foods
func (f *InMemoryFoodsStore) GetFoods() ([]Food, error) {
	foods := make([]Food, 0)
	foods = append(foods, Food{Name: "production 1", Calories: 666})
	return foods, nil
}

func main() {
	http.ListenAndServe(":5000", &FoodsServer{store: &InMemoryFoodsStore{}})
}
