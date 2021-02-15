package food

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrMissingParam error struct for displaying missing param error with specified param
type ErrMissingParam string

func (e *ErrMissingParam) Error() string {
	return string("Missing parameter: " + *e)
}

// ErrInternalServer constant for error message
const ErrInternalServer = "Internal server error"

// Food struct type
type Food struct {
	Name     string
	Calories int
}

// FoodsServer struct to use FoodsStore
type FoodsServer struct {
	Store FoodsStore
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
	foods, err := f.Store.GetFoods()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, ErrInternalServer)
	} else {
		respondWithSuccess(w, http.StatusOK, foods)
	}
}

func handlePostFood(f *FoodsServer, w http.ResponseWriter, req *http.Request) {
	var foodParam Food
	json.NewDecoder(req.Body).Decode(&foodParam)

	if foodParam.Name == "" {
		err := ErrMissingParam("Name")
		respondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if foodParam.Calories == 0 {
		err := ErrMissingParam("Calories")
		respondWithError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	food, err := f.Store.PostFood(foodParam)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, ErrInternalServer)
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
