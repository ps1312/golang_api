package users

// InMemoryUsersStore storage
type InMemoryUsersStore struct {
	users []DatabaseModel
}

func (i *InMemoryUsersStore) save(user DatabaseModel) error {
	i.users = append(i.users, user)
	return nil
}
