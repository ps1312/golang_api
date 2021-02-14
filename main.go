package main

import (
	"net/http"
)

// ErrInternalServer constant for error message
const ErrInternalServer = "Internal server error"

// ErrMissingParam constant for error message
const ErrMissingParam = "Missing parameter"

func main() {
	http.ListenAndServe(":5000", &FoodsServer{store: &InMemoryFoodsStore{}})
}
