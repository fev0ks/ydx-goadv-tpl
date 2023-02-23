package storage

import "github.com/fev0ks/ydx-goadv-tpl/internal/model"

type UserStorage interface {
	GetUser(username string) *model.User
}
