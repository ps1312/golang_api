package food

import (
	"reflect"
	"testing"
)

func TestInMemoryFoodStore(t *testing.T) {
	t.Run("Delivers empty slice of foods on empty store", func(t *testing.T) {
		store := InMemoryFoodsStore{[]Food{}}
		assertFoods(t, store, []Food{})
	})

	t.Run("Delivers slice of foods with inserted food", func(t *testing.T) {
		food := Food{"food", 1234}
		food2 := Food{"food 2", 4321}
		store := InMemoryFoodsStore{}

		store.PostFood(food)
		assertFoods(t, store, []Food{food})

		store.PostFood(food2)
		assertFoods(t, store, []Food{food, food2})
	})
}

func assertFoods(t *testing.T, store InMemoryFoodsStore, want []Food) {
	t.Helper()
	got, _ := store.GetFoods()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
