package main

import (
	"fmt"
	"net/http"
)

// Food struct type
type Food struct {
	name     string
	calories int
}

// FoodsStore interface for Food storage operations
type FoodsStore interface {
	GetFoods() []Food
}

// FoodsServer struct to use FoodsStore
type FoodsServer struct {
	store FoodsStore
}

// FoodServer handles requests for foods
func (f *FoodsServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, f.store.GetFoods())
}

func main() {
	http.ListenAndServe(":5000", &FoodsServer{})
}
