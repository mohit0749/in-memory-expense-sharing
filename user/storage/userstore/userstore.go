package userstore

import (
	"fmt"
	"sync"

	"testsplitwise/user"
)

type userStore struct {
	users sync.Map
}

func NewUserStore() *userStore {
	return &userStore{sync.Map{}}
}

func (us *userStore) AddUser(u *user.User) error {
	us.users.Store(u.GetId(), u)
	return nil
}

func (us *userStore) GetUser(id uint64) (*user.User, error) {
	value, ok := us.users.Load(id)
	if !ok {
		return nil, fmt.Errorf("user %d does not exists", id)
	}
	u, _ := value.(*user.User)
	return u, nil
}

func (us *userStore) GetAllUsers() []*user.User {
	users := make([]*user.User, 0)
	us.users.Range(func(key, value interface{}) bool {
		u, _ := value.(*user.User)
		users = append(users, u)
		return true
	})
	return users
}
