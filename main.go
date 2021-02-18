package main

import (
	"api/encryption"
	"api/food"
	"api/signer"
	"api/user"
	"net/http"
)

func main() {
	http.Handle("/foods", &food.FoodsServer{Store: &food.InMemoryFoodsStore{Foods: []food.Food{}}})
	http.Handle("/users", &user.Server{
		Encrypter: &encryption.BCryptEncrypter{},
		Store:     &user.InMemoryUsersStore{Users: []user.DatabaseModel{}},
		Signer:    &signer.JWTSigner{Key: []byte("any secret key")},
	})
	http.ListenAndServe(":5000", nil)
}
