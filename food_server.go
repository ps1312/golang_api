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
	} else {
		var foodParam Food
		json.NewDecoder(req.Body).Decode(&foodParam)
		f.store.PostFood(foodParam)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, ErrInternalServer)
	}
}
