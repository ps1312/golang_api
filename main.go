package main

import (
	"api/food"
	"net/http"
)

func main() {
	http.ListenAndServe(":5000", &food.FoodsServer{Store: &food.InMemoryFoodsStore{Foods: []food.Food{}}})
}
