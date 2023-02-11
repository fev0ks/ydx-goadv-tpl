package service

import "github.com/fev0ks/ydx-goadv-tpl/model"

type UserService interface {
	Register(auth *model.Auth) error
	Login(auth *model.Auth) error
}

type userService struct {
}

func NewUserService() UserService {
	return &userService{}
}

func (us *userService) Register(auth *model.Auth) error {
	return nil
}

func (us *userService) Login(auth *model.Auth) error {
	return nil
}
