package repository

import (
	"GoLangProject/webhook/internal/domain"
	"GoLangProject/webhook/internal/repository/dao"
	"context"
	"time"
)

var (
	DuplicatedError    = dao.DuplicatedError
	ErrUserNotFound    = dao.ErrRecordNotFound
	ErrProfileNotFound = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) CreateUser(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		ID:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		Ctime:    time.Unix(u.Ctime/1000, u.Ctime%1000*1000000),
	}
}

func (repo *UserRepository) ModifyProfile(ctx context.Context, u domain.UserProfile) error {
	birthdayUnixMilli := u.Birthday.UnixMilli()
	return repo.dao.ModifyProfile(ctx, dao.UserProfile{
		ID:       u.ID,
		Nickname: u.Nickname,
		Birthday: birthdayUnixMilli,
		AboutMe:  u.AboutMe,
	})
}

func (repo *UserRepository) FindByUserID(ctx context.Context, userID int64) (domain.UserProfile, error) {
	daoUserProfile, err := repo.dao.FindByUserID(ctx, userID)
	if err != nil {
		return domain.UserProfile{}, err
	}
	return repo.profileToDomain(daoUserProfile), err
}

func (repo *UserRepository) profileToDomain(profile dao.UserProfile) domain.UserProfile {
	return domain.UserProfile{
		Nickname: profile.Nickname,
		Birthday: time.Unix(profile.Birthday/1000, (profile.Birthday%1000)*1000000),
		AboutMe:  profile.AboutMe,
	}
}
