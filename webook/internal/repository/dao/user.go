package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

type UserDAO struct {
	db *gorm.DB
}

type User struct {
	Id       int64  `gorm:"primaryKey, autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	// 时区问题 UTC0 的毫秒数
	Ctime int64
	Utime int64
}

type UserProfile struct {
	ID       int64 `gorm:"primaryKey"`
	Nickname string
	Birthday int64
	AboutMe  string
}

var (
	DuplicatedError   = errors.New("邮箱冲突错误，请重新选择邮箱")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

// Insert 存储直接函数 用于直接存储相关信息
func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		if me.Number == uint16(1062) {
			return DuplicatedError
		}
	}
	return err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	if errors.Is(err, ErrRecordNotFound) {
		return User{}, ErrRecordNotFound
	}
	return u, err
}

func (dao *UserDAO) ModifyProfile(ctx context.Context, up UserProfile) error {
	// Put Data Into DataBase
	err := dao.db.WithContext(ctx).Save(&up).Error
	if err != nil {

		return err
	}
	return err
}

func (dao *UserDAO) FindByUserID(ctx context.Context, userID int64) (UserProfile, error) {
	var userProfile UserProfile
	err := dao.db.WithContext(ctx).Where("ID=?", userID).First(&userProfile).Error
	if errors.Is(err, ErrRecordNotFound) {
		return UserProfile{}, err
	}
	return userProfile, err
}
