package main

// FoodsStore interface for Food storage operations
type FoodsStore interface {
	GetFoods() ([]Food, error)
	PostFood(food Food) (Food, error)
}

// InMemoryFoodsStore in memory store for testing
type InMemoryFoodsStore struct{}

// GetFoods returns foods
func (f *InMemoryFoodsStore) GetFoods() ([]Food, error) {
	foods := make([]Food, 0)
	return foods, nil
}

// PostFood saves food
func (f *InMemoryFoodsStore) PostFood(food Food) (Food, error) {
	return Food{}, nil
}
