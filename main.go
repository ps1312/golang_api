package main

import (
	"net/http"
)

// ErrMissingParam error struct for displaying missing param error with specified param
type ErrMissingParam string

func (e *ErrMissingParam) Error() string {
	return string("Missing parameter: " + *e)
}

func (e *ErrMissingParam) String() string {
	return string("Missing parameter: " + *e)
}

// ErrInternalServer constant for error message
const ErrInternalServer = "Internal server error"

func main() {
	http.ListenAndServe(":5000", &FoodsServer{store: &InMemoryFoodsStore{[]Food{}}})
}
