package users

import "testing"

func Test(t *testing.T) {
	model := DatabaseModel{Name: "any-name", Email: "any@mail.com", password: "any-password"}
	store := InMemoryUsersStore{}

	store.save(model)

	if len(store.Users) != 1 {
		t.Errorf("got %d, want 1", len(store.Users))
	}
}
