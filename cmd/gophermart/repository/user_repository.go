package repository

import (
	"context"
	"fmt"
	"github.com/fev0ks/ydx-goadv-tpl/model"
	"github.com/jackc/pgx/v4"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, username string) (*model.User, error)
}

type userRepository struct {
	db DBProvider
}

func NewUserRepository(db DBProvider) UserRepository {
	return &userRepository{db}
}

func (u userRepository) CreateUser(ctx context.Context, user *model.User) error {
	_, err := u.db.GetConnection().Exec(ctx, "insert into users(username, password) values($1, $2)", user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (u userRepository) GetUser(ctx context.Context, username string) (*model.User, error) {
	result := u.db.GetConnection().QueryRow(ctx, "select username, password from users where username=$1", username)
	user := &model.User{}
	err := result.Scan(&user.Username, &user.Password)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get stored creds of '%s': %v", username, err)
	}
	return user, nil
}
