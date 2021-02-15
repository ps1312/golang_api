package main

import (
	"reflect"
	"testing"
)

func Test(t *testing.T) {
	t.Run("Delivers empty slice of foods on empty store", func(t *testing.T) {
		store := InMemoryFoodsStore{}
		got, _ := store.GetFoods()
		want := []Food{}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
