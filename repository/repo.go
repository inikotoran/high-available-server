package repository

import (
	"github.com/inikotoran/high-available-server/model"
)

type Repo interface {
	Save(model.User) error
	Get(username string) (*model.User, error)
}
