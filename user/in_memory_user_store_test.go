package users

import "testing"

type InMemoryUsersStore struct {
	users []DatabaseModel
}

func (i *InMemoryUsersStore) save(user DatabaseModel) error {
	i.users = append(i.users, user)
	return nil
}

func Test(t *testing.T) {
	model := DatabaseModel{Name: "any-name", Email: "any@mail.com", password: "any-password"}
	store := InMemoryUsersStore{}

	store.save(model)

	if len(store.users) != 1 {
		t.Errorf("got %d, want 1", len(store.users))
	}
}
