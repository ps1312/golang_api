package user

import (
	"reflect"
	"testing"
)

func TestInMemoryStore(t *testing.T) {
	t.Run("Saves user on storage correctly", func(t *testing.T) {
		model := DatabaseModel{Name: "any-name", Email: "any@mail.com", password: "any-password"}
		store := InMemoryUsersStore{}

		store.save(model)

		if len(store.Users) != 1 {
			t.Errorf("got %d, want 1", len(store.Users))
		}
	})

	t.Run("Delivers saved user on getAll", func(t *testing.T) {
		model := DatabaseModel{Name: "any-name", Email: "any@mail.com", password: "any-password"}
		store := InMemoryUsersStore{}

		store.save(model)

		got, _ := store.getAll()
		want := []DatabaseModel{model}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
