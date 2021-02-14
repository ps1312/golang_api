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
		respondWithError(w, http.StatusInternalServerError, ErrInternalServer)
	} else {
		respondWithSuccess(w, http.StatusOK, foods)
	}
}

func handlePostFood(f *FoodsServer, w http.ResponseWriter, req *http.Request) {
	var foodParam Food
	json.NewDecoder(req.Body).Decode(&foodParam)
	food, err := f.store.PostFood(foodParam)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, ErrInternalServer)
	} else if foodParam.Calories == 0 || foodParam.Name == "" {
		respondWithError(w, http.StatusUnprocessableEntity, ErrMissingParam)
	} else {
		respondWithSuccess(w, http.StatusCreated, food)
	}
}

func respondWithError(w http.ResponseWriter, status int, err string) {
	w.WriteHeader(status)
	fmt.Fprint(w, err)
}

func respondWithSuccess(w http.ResponseWriter, status int, body interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
