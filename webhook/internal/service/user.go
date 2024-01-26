package service

import (
	"GoLangProject/webhook/internal/domain"
	"GoLangProject/webhook/internal/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var DuplicatedError = repository.DuplicatedError

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, user domain.User) error {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("not Generated Password into Encrypted Code")
	}
	user.Password = string(encrypted)
	return svc.repo.CreateUser(ctx, user)
}

func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	// 检查密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}
