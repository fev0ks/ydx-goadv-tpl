package service

import (
	"context"
	"github.com/fev0ks/ydx-goadv-tpl/model"
	"github.com/fev0ks/ydx-goadv-tpl/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, userRequest *model.UserRequest) error
	GetUser(ctx context.Context, username string) (*model.User, error)
	IsCorrectUserPassword(ctx context.Context, userRequest *model.UserRequest) (bool, error)
	IsExistUsername(ctx context.Context, username string) (bool, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (us *userService) CreateUser(ctx context.Context, userRequest *model.UserRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), 8)
	if err != nil {
		return err
	}
	newUser := &model.User{
		Username: userRequest.Username,
		Password: hashedPassword,
	}
	return us.userRepo.CreateUser(ctx, newUser)
}

func (us *userService) GetUser(ctx context.Context, username string) (*model.User, error) {
	return us.userRepo.GetUser(ctx, username)
}

func (us *userService) IsCorrectUserPassword(ctx context.Context, userRequest *model.UserRequest) (bool, error) {
	user, err := us.GetUser(ctx, userRequest.Username)
	if user == nil || err != nil {
		return false, err
	}
	if err = bcrypt.CompareHashAndPassword(user.Password, []byte(userRequest.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (us *userService) IsExistUsername(ctx context.Context, username string) (bool, error) {
	user, err := us.userRepo.GetUser(ctx, username)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}
