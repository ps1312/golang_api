package user

// InMemoryUsersStore storage
type InMemoryUsersStore struct {
	Users []DatabaseModel
}

func (i *InMemoryUsersStore) save(user DatabaseModel) error {
	i.Users = append(i.Users, user)
	return nil
}

func (i *InMemoryUsersStore) getAll() ([]DatabaseModel, error) {
	return i.Users, nil
}
