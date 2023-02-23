package service

import (
	"context"
	"github.com/fev0ks/ydx-goadv-tpl/internal/model"
	"github.com/fev0ks/ydx-goadv-tpl/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, userRequest *model.UserRequest) (*model.User, error)
	GetUser(ctx context.Context, username string) (*model.User, error)
	ValidatePassword(ctx context.Context, user *model.User, password string) (bool, error)
	IsExistUsername(ctx context.Context, username string) (bool, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (us *userService) CreateUser(ctx context.Context, userRequest *model.UserRequest) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), 8)
	if err != nil {
		return nil, err
	}
	newUser := &model.User{
		Username: userRequest.Login,
		Password: hashedPassword,
	}
	return us.userRepo.CreateUser(ctx, newUser)
}

func (us *userService) GetUser(ctx context.Context, username string) (*model.User, error) {
	return us.userRepo.GetUser(ctx, username)
}

func (us *userService) ValidatePassword(_ context.Context, user *model.User, password string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
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
