package main

import (
	"net/http"
)

// ErrInternalServer constant for error message
const ErrInternalServer = "Internal server error"

func main() {
	http.ListenAndServe(":5000", &FoodsServer{store: &InMemoryFoodsStore{}})
}
