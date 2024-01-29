package service

import (
	"GoLangProject/webhook/internal/domain"
	"GoLangProject/webhook/internal/repository"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	DuplicatedError          = repository.DuplicatedError
	ErrInvalidUserOrPassword = repository.ErrUserNotFound
	ErrUserProfileNotFound   = repository.ErrProfileNotFound
)

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
	// Check Email is correct
	u, err := svc.repo.FindByEmail(ctx, email)
	if errors.Is(err, ErrInvalidUserOrPassword) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// Check Password is correct
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (svc *UserService) Edit(ctx *gin.Context, userID int64, nickname string, birthday time.Time, aboutme string) error {
	err := svc.repo.ModifyProfile(ctx, domain.UserProfile{
		ID:       userID,
		Nickname: nickname,
		Birthday: birthday,
		AboutMe:  aboutme,
	})
	if err != nil {
		return err
	}
	return err
}

func (svc *UserService) Profile(ctx *gin.Context, userId int64) (domain.UserProfile, error) {
	domainUserProfile, err := svc.repo.FindByUserID(ctx, userId)
	return domainUserProfile, err
}
