package repository

import (
	"context"
	"fmt"
	"github.com/fev0ks/ydx-goadv-tpl/model"
	"github.com/jackc/pgx/v4"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUser(ctx context.Context, username string) (*model.User, error)
}

type userRepository struct {
	db DBProvider
}

func NewUserRepository(db DBProvider) UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	var userID int
	tx, err := ur.db.GetConnection().Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return nil, err
	}
	row := tx.QueryRow(ctx, "insert into users(username, password) values($1, $2) RETURNING user_id", user.Username, user.Password)
	err = row.Scan(&userID)
	if err != nil {
		return nil, err
	}
	user.UserID = userID
	_, err = tx.Exec(ctx, "insert into user_balance(user_id, current, withdraw) values($1, $2, $3)", userID, 0, 0)
	if err != nil {
		return nil, err
	}
	tx.Commit(ctx)
	return user, nil
}

func (ur *userRepository) GetUser(ctx context.Context, username string) (*model.User, error) {
	result := ur.db.GetConnection().QueryRow(ctx, "select user_id, username, password from users where username=$1", username)
	user := &model.User{}
	err := result.Scan(&user.UserID, &user.Username, &user.Password)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get stored creds of '%s': %v", username, err)
	}
	return user, nil
}
