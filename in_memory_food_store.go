package main

// FoodsStore interface for Food storage operations
type FoodsStore interface {
	GetFoods() ([]Food, error)
	PostFood(food Food) (Food, error)
}

// InMemoryFoodsStore in memory store for testing
type InMemoryFoodsStore struct {
	foods []Food
}

// GetFoods returns foods
func (f *InMemoryFoodsStore) GetFoods() ([]Food, error) {
	return f.foods, nil
}

// PostFood saves food
func (f *InMemoryFoodsStore) PostFood(food Food) (Food, error) {
	f.foods = append(f.foods, food)
	return food, nil
}
