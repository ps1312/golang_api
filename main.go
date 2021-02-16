package main

import (
	"api/food"
	"api/user"
	"net/http"
)

func main() {
	http.Handle("/foods", &food.FoodsServer{Store: &food.InMemoryFoodsStore{Foods: []food.Food{}}})
	http.Handle("/users", &user.Server{Encrypter: &user.BCryptEncrypter{}, Store: &user.InMemoryUsersStore{Users: []user.DatabaseModel{}}})
	http.ListenAndServe(":5000", nil)
}
