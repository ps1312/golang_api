package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Food struct type
type Food struct {
	Name     string
	Calories int
}

// FoodsServer struct to use FoodsStore
type FoodsServer struct {
	store FoodsStore
}

// FoodServer handles requests for foods
func (f *FoodsServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		handleGetFoods(f, w, req)
	} else {
		handlePostFood(f, w, req)
	}
}

func handleGetFoods(f *FoodsServer, w http.ResponseWriter, req *http.Request) {
	foods, err := f.store.GetFoods()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ErrInternalServer)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(foods)
	}
}

func handlePostFood(f *FoodsServer, w http.ResponseWriter, req *http.Request) {
	var foodParam Food
	json.NewDecoder(req.Body).Decode(&foodParam)
	food, err := f.store.PostFood(foodParam)

	if foodParam.Calories == 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, ErrMissingParam)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ErrInternalServer)
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(food)
	}

}
