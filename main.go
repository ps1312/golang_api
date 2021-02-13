package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
	foods, err := f.store.GetFoods()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(foods)
	}
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
