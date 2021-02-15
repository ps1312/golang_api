package main

import (
	"reflect"
	"testing"
)

func TestInMemoryGetFoods(t *testing.T) {
	t.Run("Delivers empty slice of foods on empty store", func(t *testing.T) {
		store := InMemoryFoodsStore{}
		got, _ := store.GetFoods()
		want := 0

		if len(got) != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Delivers slice of stored foods", func(t *testing.T) {
		food := Food{"food", 1234}
		store := InMemoryFoodsStore{}

		store.PostFood(food)

		got, _ := store.GetFoods()
		want := []Food{food}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
