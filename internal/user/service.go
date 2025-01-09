package user

import (
	"catalog-digital-product/internal/custom"
	"catalog-digital-product/internal/helper"
	"catalog-digital-product/internal/token"
	"context"
	"database/sql"
)

type UserService interface {
	Login(ctx context.Context, input LoginInputUser) (User, string, error)
	UpdatePassword(ctx context.Context, input UpdatePasswordInputUser, currentUserId int) error
}

type UserServiceImpl struct {
	DB             *sql.DB
	UserRepository UserRepository
	TokenService   token.TokenService
}

func NewUserService(DB *sql.DB, userRepository UserRepository, tokenService token.TokenService) UserService {
	return &UserServiceImpl{DB: DB, UserRepository: userRepository, TokenService: tokenService}
}

func (u *UserServiceImpl) Login(ctx context.Context, input LoginInputUser) (User, string, error) {
	tx, err := u.DB.Begin()
	if err != nil {
		return User{}, "", err
	}
	defer helper.HandleTransaction(tx, &err)

	user, err := u.UserRepository.FindByUsername(ctx, tx, input.Username)
	if err != nil {
		return User{}, "", err
	}
	if user.Id <= 0 {
		return User{}, "", custom.ErrNotFound
	}
	isSame, err := helper.ComparePassword(user.Password, input.Password)
	if err != nil {
		return User{}, "", err
	}
	if !isSame {
		return User{}, "", custom.ErrUnauthorized
	}

	token, err := u.TokenService.GenerateToken(user.Id)
	if err != nil {
		return User{}, "", err
	}

	return user, token, nil
}

func (u *UserServiceImpl) UpdatePassword(ctx context.Context, input UpdatePasswordInputUser, currentUserId int) error {
	tx, err := u.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.HandleTransaction(tx, &err)

	user, err := u.UserRepository.FindById(ctx, tx, currentUserId)
	if err != nil {
		return err
	}
	if user.Id <= 0 {
		return custom.ErrNotFound
	}

	user.Password = input.Password
	err = u.UserRepository.Update(ctx, tx, user)
	if err != nil {
		return err
	}

	return nil
}
