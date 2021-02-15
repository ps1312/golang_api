package food

// FoodsStore interface for Food storage operations
type FoodsStore interface {
	GetFoods() ([]Food, error)
	PostFood(food Food) (Food, error)
}

// InMemoryFoodsStore in memory store for testing
type InMemoryFoodsStore struct {
	Foods []Food
}

// GetFoods returns Foods
func (f *InMemoryFoodsStore) GetFoods() ([]Food, error) {
	return f.Foods, nil
}

// PostFood saves food
func (f *InMemoryFoodsStore) PostFood(food Food) (Food, error) {
	f.Foods = append(f.Foods, food)
	return food, nil
}
