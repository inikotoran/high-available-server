package repository

import (
	"fmt"
	"github.com/inikotoran/high-available-server/model"
	"sync"
)

type InMemoryRepo struct {
	users     map[string]model.User
	usersLock sync.Mutex
}

func (i *InMemoryRepo) Save(user model.User) error {
	i.usersLock.Lock()
	defer i.usersLock.Unlock()

	i.users[user.Username] = user
	return nil
}

func (i *InMemoryRepo) Get(username string) (*model.User, error) {
	i.usersLock.Lock()
	user, found := i.users[username]
	i.usersLock.Unlock()
	if !found {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}

func NewInMemoryRepo() Repo {
	return &InMemoryRepo{
		users: make(map[string]model.User),
	}
}
